package main

import (
	"crypto/tls"
	"dynamicpath/lib/MonitorInfo"
	"dynamicpath/lib/path_util"
	"dynamicpath/src/monitor/http_server"
	"dynamicpath/src/monitor/monitor_context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var host string
var bandwidthFormat = regexp.MustCompile(`(v=(\w+),)?rx=(\d+.\d*),tx=(\d+.\d*)`)

func parseInterfaceLinkInfo() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return
	}
	ifaceBandwidth := make([]string, len(ifaces))
	for i, iface := range ifaces {
		flag.StringVar(&ifaceBandwidth[i], iface.Name, "", "UPF Connnection IP in this interface.")
	}
	flag.Parse()
	for i, bandwidth := range ifaceBandwidth {
		if bandwidth != "" {
			info := MonitorInfo.BandwidthInfo{}
			matches := bandwidthFormat.FindStringSubmatch(bandwidth)
			if matches != nil {
				if matches[2] != "" {
					// Has Virtual Interface
					monitor_context.VirtualIface[ifaces[i].Name] = &matches[2]
				}
				info.Rx, _ = strconv.ParseFloat(matches[3], 64)
				info.Tx, _ = strconv.ParseFloat(matches[4], 64)
			} else {
				log.Panicf("Error Interface format: %s, Usage: -eth rx=30,tx=30 (Unit is mbps)", bandwidth)
			}
			monitor_context.IfacesBandwidth[ifaces[i].Name] = info
		}
	}
}

func main() {
	// runtime.GOMAXPROCS(10)

	flag.StringVar(&host, "h", ":9487", "set `host` value")
	parseInterfaceLinkInfo()
	flag.Parse()

	monitorLogPath := path_util.ModulePath("dynamicpath/support/TLS/monitor_" + host + ".log")
	mointorPemPath := path_util.ModulePath("dynamicpath/support/TLS/monitor.pem")
	mointorKeyPath := path_util.ModulePath("dynamicpath/support/TLS/monitor.key")

	router := http_server.NewRouter()
	fmt.Printf("Server Listen on: %s\n", host)

	keylogFile, err := os.OpenFile(monitorLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic(err)
	}
	server := &http.Server{
		Addr: host,
		TLSConfig: &tls.Config{
			KeyLogWriter: keylogFile,
		},
		Handler: router,
	}

	err = server.ListenAndServeTLS(mointorPemPath, mointorKeyPath)
	if err != nil {
		log.Panic(err)
	}
}
