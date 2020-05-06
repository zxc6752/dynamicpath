package http_server

import (
	"dynamicpath/lib/MonitorInfo"
	"dynamicpath/src/monitor/monitor_context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"net/http"
)

var sem = semaphore.NewWeighted(1)
var nSec int = 2

func GetMonitorData(c *gin.Context) {
	if monitor_context.UpfMonitorData.PacketRates == nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[Error] Monitor has not started yet")
		return
	} else if !sem.TryAcquire(1) {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[WARN] Monitor has not finished yet")
		return
	}
	defer sem.Release(1)
	monitor_context.UpdateMonitorData()
	c.JSON(http.StatusOK, monitor_context.UpfMonitorData)
	monitor_context.StartCollectData()
}

func StartMonitor(c *gin.Context) {
	if monitor_context.UpfMonitorData.PacketRates != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[Error] Monitor has started")
		return
	} else if !sem.TryAcquire(1) {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[WARN] Monitor has started")
		return
	}
	defer sem.Release(1)
	var monitorStart MonitorInfo.MonitorStartRequest
	if err := c.ShouldBindJSON(&monitorStart); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	monitor_context.InitInterfacesInfo(monitorStart.LinkIPInfo)
	// for name, iface := range monitor_context.IfacesLinkInfo {
	// 	fmt.Println(name, iface.Link, iface.Speed)
	// }
	monitor_context.InitUpfMonitorInfo()
	c.JSON(http.StatusOK, monitor_context.IfacesInfo)
	monitor_context.InitStatistics()
	monitor_context.RefreshPacketRate(2)
	monitor_context.StartCollectData()
}

// func timerStart(period int) {
// 	timeDuraion := time.Duration(period-monitor_context.NSec)*time.Second - time.Duration(monitor_context.ProcessDelay)*time.Millisecond
// 	monitor_context.Timer = time.NewTimer(timeDuraion)

// 	go func() {
// 		select {
// 		case <-monitor_context.Timer.C:
// 			sem.Acquire(context.Background(), 1)
// 			sem.Release(1)
// 		case <-time.After(timeDuraion + 100*time.Millisecond):
// 			fmt.Println("update timer changed")
// 		}
// 	}()

// }
