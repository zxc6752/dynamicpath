package monitor_context

import (
	"dynamicpath/lib/MonitorInfo"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type statisticsData struct {
	Timestamp time.Time
	RxPackets uint64
	TxPackets uint64
}

// const NSec int = 10

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var VirtualIface map[string]*string // physical Eth -> Virtual ifb , Physical-RX=Virtual-TX this is for ingress bandwidth limit
var IfacesBandwidth map[string]MonitorInfo.BandwidthInfo
var IfacesInfo MonitorInfo.InterfacesInfo
var UpfMonitorData MonitorInfo.UPFMonitorData
var Timer *time.Timer

var remoteIpPing map[string]*exec.Cmd
var UpdatePacketRateMtx sync.Mutex
var pingMtx sync.Mutex
var packetRateMtx sync.Mutex
var statisticsDatas map[string]*statisticsData
var wg sync.WaitGroup
var wgPacketRate sync.WaitGroup
var pingFormat = regexp.MustCompile(`received, (\d+\.?\d*)% .*\n.* = \d+.\d*/(\d+.\d*)/`)

func init() {
	VirtualIface = make(map[string]*string)
	IfacesBandwidth = make(map[string]MonitorInfo.BandwidthInfo)
	IfacesInfo.BandwidthMaps = make(map[string]MonitorInfo.BandwidthInfo)
	IfacesInfo.RemoteIpMaps = make(map[string]string)
	statisticsDatas = make(map[string]*statisticsData)
	remoteIpPing = make(map[string]*exec.Cmd)
}

func InitUpfMonitorInfo() {
	UpfMonitorData.PacketRates = make(map[string]MonitorInfo.BandwidthInfo)
	UpfMonitorData.ConnectionInfos = make(map[string]*MonitorInfo.ConnectionInfo)
	for ip, _ := range remoteIpPing {
		UpfMonitorData.ConnectionInfos[ip] = &MonitorInfo.ConnectionInfo{}
	}
}

func InitInterfacesInfo(linkIPInfo []MonitorInfo.MonitorIPInfo) {
	for _, ipInfo := range linkIPInfo {
		cmd := fmt.Sprintf("ip route get \"%s\" | grep -Po '(?<=(dev )).*(?= src| proto)'", ipInfo.IP)
		output, err := exec.Command("bash", "-c", cmd).Output()
		check(err)
		var ifaceName string
		_, err = fmt.Sscanf(string(output), "%s", &ifaceName)
		check(err)
		if _, exist := IfacesInfo.BandwidthMaps[ifaceName]; !exist {
			if info, ok := IfacesBandwidth[ifaceName]; ok {
				IfacesInfo.BandwidthMaps[ifaceName] = info
			} else {
				IfacesInfo.BandwidthMaps[ifaceName] = MonitorInfo.BandwidthInfo{
					Rx: 1000.0,
					Tx: 1000.0,
				}
				log.Printf("Interface[%s] Rx and Tx Bandwidth not defined, use 1Gbps as RX/TX", ifaceName)
			}
		}
		IfacesInfo.RemoteIpMaps[ipInfo.IP] = ifaceName
		if ipInfo.Update {
			remoteIpPing[ipInfo.IP] = &exec.Cmd{}
		}
	}
}

func InitStatistics() {
	UpdatePacketRateMtx.Lock()
	wgPacketRate.Add(len(IfacesInfo.BandwidthMaps))
	mtx := sync.Mutex{}
	for ifaceName, _ := range IfacesInfo.BandwidthMaps {
		go func(name string) {
			statsData := statisticsData{}
			statsData.RxPackets, statsData.TxPackets = getCurrStatistics(name, VirtualIface[name])
			statsData.Timestamp = time.Now()
			mtx.Lock()
			statisticsDatas[name] = &statsData
			mtx.Unlock()
			wgPacketRate.Done()
		}(ifaceName)
	}
	wgPacketRate.Wait()
}

// func updateCpuUsage() {
// 	stdout, err := exec.Command("bash", "-c", "top -b -n2 | grep 'Cpu(s)' | awk '{print $2+$4}' | tail -n 1").Output()
// 	check(err)
// 	_, err = fmt.Sscanf(string(stdout), "%f", &upfMonitotData.CpuUsage)
// 	check(err)
// 	wg.Done()
// }

func doPing(linkIP string, conn *MonitorInfo.ConnectionInfo) {
	defer wg.Done()
	{
		pingMtx.Lock()
		remoteIpPing[linkIP] = exec.Command("ping", linkIP, "-I", IfacesInfo.RemoteIpMaps[linkIP], "-i", "0.1", "-q")
		pingMtx.Unlock()
		output, _ := remoteIpPing[linkIP].Output()
		pingMtx.Lock()
		// remoteIpPing[linkIP].Process.Kill()
		remoteIpPing[linkIP] = nil
		pingMtx.Unlock()
		matches := pingFormat.FindStringSubmatch(string(output))
		if matches != nil {
			loss, _ := strconv.ParseFloat(matches[1], 64)
			conn.PacketLoss = loss / 100.0
			delay, _ := strconv.ParseFloat(matches[2], 64)
			conn.DelayTime = delay
		} else {
			log.Printf("ping %s failed[%s]", linkIP, output)
			conn.PacketLoss = 1.0
			conn.DelayTime = 10000.0 // 10s delay
			return
		}
	}
}

func updateConnectionInfo() {
	for linkIP, conn := range UpfMonitorData.ConnectionInfos {
		go doPing(linkIP, conn)
	}
}

func getCurrStatistics(ifaceName string, virtualIface *string) (rx, tx uint64) {
	var rxBytes, txBytes []byte
	var err error
	if virtualIface != nil {
		rxBytes, err = exec.Command("cat", "/sys/class/net/"+*virtualIface+"/statistics/tx_bytes").Output()
		check(err)
	} else {
		rxBytes, err = exec.Command("cat", "/sys/class/net/"+ifaceName+"/statistics/rx_bytes").Output()
		check(err)
	}
	txBytes, err = exec.Command("cat", "/sys/class/net/"+ifaceName+"/statistics/tx_bytes").Output()
	check(err)
	_, err = fmt.Sscanf(string(rxBytes), "%d\n", &rx)
	check(err)
	_, err = fmt.Sscanf(string(txBytes), "%d\n", &tx)
	check(err)
	return
}

func updatePacketRate(ifaceName string) {
	defer wgPacketRate.Done()
	currRx, currTx := getCurrStatistics(ifaceName, VirtualIface[ifaceName])
	statsData := statisticsDatas[ifaceName]
	nSeconds := time.Since(statsData.Timestamp).Seconds()
	info := MonitorInfo.BandwidthInfo{}
	info.Rx = float64(currRx-statsData.RxPackets) / (math.Pow(2, 20) * nSeconds) * 8
	info.Tx = float64(currTx-statsData.TxPackets) / (math.Pow(2, 20) * nSeconds) * 8
	statsData.Timestamp = time.Now()
	statsData.RxPackets, statsData.TxPackets = currRx, currTx
	packetRateMtx.Lock()
	UpfMonitorData.PacketRates[ifaceName] = info
	packetRateMtx.Unlock()

}

func RefreshPacketRate(nSec int) {
	duration := time.Millisecond * time.Duration(1000.0*nSec)
	timer := time.NewTimer(duration)

	go func() {
		select {
		case <-timer.C:
			UpdatePacketRateMtx.Lock()
			wgPacketRate.Add(len(IfacesInfo.BandwidthMaps))
			for ifaceName, _ := range IfacesInfo.BandwidthMaps {
				go updatePacketRate(ifaceName)
			}
			wgPacketRate.Wait()
			UpdatePacketRateMtx.Unlock()
			RefreshPacketRate(nSec)
		case <-time.After(duration + 100*time.Millisecond):
		}
	}()
	return
}

func UpdateMonitorData() {
	UpdatePacketRateMtx.Lock()
	pingMtx.Lock()
	for _, cmd := range remoteIpPing {
		if cmd != nil {
			cmd.Process.Signal(syscall.SIGINT)
		}
	}
	pingMtx.Unlock()
	wg.Wait()
}

func StartCollectData() {
	// setStatistics()
	UpdatePacketRateMtx.Unlock()
	wg.Add(len(UpfMonitorData.ConnectionInfos))
	updateConnectionInfo()
}
