package lb_util

import (
	"dynamicpath/lib/up_topology"
	"dynamicpath/src/load_balancer/lb_client"
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/logger"
	"sync"
	"time"
)

type RemainRateOperation string

const (
	AddRemainRate         RemainRateOperation = "Add"
	MinusRemainRate       RemainRateOperation = "Minus"
	LoadBalancerType_RIRM int                 = 0
	LoadBalancerType_DP   int                 = 3
)

var DP_ShortestPathThreshlod int = -1
var AvailableBandwidth []float64
var mtx sync.Mutex // Mutex for Query Data and Update BestPath

func UpdateMonitorData(nSec int) *time.Timer {
	timer := time.NewTimer(time.Second * time.Duration(nSec))

	go func() {
		select {
		case <-timer.C:
			lb_context.Wg.Add(1)
			mtx.Lock()
			self := lb_context.LB_Self()
			if self.Topology.RefreshTimer != nil {
				self.Topology.RefreshTimer.Stop()
				self.Topology.RefreshTimer = nil
			}
			// Wait for Pdu Session Establishment
			lb_context.Mtx.Lock()
			lb_context.Mtx.Unlock()
			queryMonitorData(nSec)
			RefreshAll(nSec)
			mtx.Unlock()
			lb_context.Wg.Done()
		case <-time.After(time.Second*time.Duration(nSec) + 100*time.Millisecond):
			logger.UtilLog.Debugf("timer closed")
		}
	}()
	return timer
}

func findDpThresh(pathList []up_topology.Path) int {
	result := len(pathList)
	costRecur := *pathList[0].Cost
	for i, path := range pathList {
		if *path.Cost > costRecur {
			result = i
			break
		}
		availBandwidth := *path.RemainRate * 0.8
		AvailableBandwidth = append(AvailableBandwidth, availBandwidth)
	}
	return result
}

func RefreshAll(nSec int) {
	topology := lb_context.LB_Self().Topology
	UpdatePathAll(topology.PathListAll)
	if lb_context.LB_Self().LoadBalancerType == LoadBalancerType_DP {
		if DP_ShortestPathThreshlod == -1 { // shortest path mode
			PathMergeSort(topology.PathListAll)
			DP_ShortestPathThreshlod = findDpThresh(topology.PathListAll)
		} else if DP_ShortestPathThreshlod == -2 { // load balancing mode
			PathMergeSort(topology.PathListAll)
		}
	} else {
		PathMergeSort(topology.PathListAll)
		// topology.PathListWorst, topology.RemainRateWorst = GetSubPathList(topology.PathListAll[topology.HalfLen:], WorstPath)
		// logger.UtilLog.Debugf("MeanWorstRate: %.2f\n", topology.RemainRateWorst)
		// for i, path := range topology.PathListWorst {
		// 	fmt.Printf("\n%d: %v, Cost-%.2f, Rate=%.2f\n", i, path.NodeIdList, *path.Cost, *path.RemainRate)
		// }
		topology.PathThresh = FindPathThreshold(topology.PathListAll)
		if topology.PathThresh == -1 {
			logger.UtilLog.Warnf("All paths are overflow")
			topology.PathListBest = nil
		} else {
			topology.PathListBest, topology.RemainRateBest = GetSubPathList(topology.PathListAll[:topology.PathThresh+1], BestPath)
			for _, path := range topology.PathListAll[:topology.PathThresh+1] {
				logger.UtilLog.Debugf("PathId: %d Load: %.2f Cost %.2f", path.Id, *path.RemainRate, *path.Cost)
			}
			logger.UtilLog.Debugf("MeanBestRate: %.2f\n", topology.RemainRateBest)
			RefreshBestPath(float64(nSec) / 10.0)
		}
	}

	lb_context.LB_Self().QueryMonitorTimer = UpdateMonitorData(nSec)

	// for i, path := range topology.PathListBest {
	// fmt.Printf("\n%d: %v, Cost-%.2f, Rate=%.2f\n", i, path.NodeIdList, *path.Cost, *path.RemainRate)
	// }
	// TODO: dynamically decide update time
	// refresh 10 times per update
	// lb_client.SendWorstPathList(topology.PathListWorst)
}

func RefreshBestPath(nSec float64) *time.Timer {
	duration := time.Millisecond * time.Duration(1000.0*nSec)
	timer := time.NewTimer(duration)

	go func() {
		select {
		case <-timer.C:
			lb_context.Wg.Add(1)
			// TODO: Semphore
			mtx.Lock()
			// Wait for Pdu Session Establishment
			lb_context.Mtx.Lock()
			lb_context.Mtx.Unlock()
			topology := lb_context.LB_Self().Topology
			UpdatePathRemainRate(topology.PathListAll[:topology.PathThresh+1])
			topology.PathListBest, topology.RemainRateBest = GetSubPathList(topology.PathListAll[:topology.PathThresh+1], BestPath)
			logger.UtilLog.Debugf("MeanBestRate: %.2f\n", topology.RemainRateBest)
			// for i, path := range topology.PathListBest {
			// 	fmt.Printf("%d: %v, Cost-%.2f, Rate=%.2f\n", i, path.NodeIdList, *path.Cost, *path.RemainRate)
			// }
			topology.RefreshTimer = RefreshBestPath(nSec)
			mtx.Unlock()
			lb_context.Wg.Done()
		case <-time.After(duration + 100*time.Millisecond):
			logger.UtilLog.Debugf("timer closed")
		}
	}()
	return timer
}

func queryMonitorData(nSec int) {
	self := lb_context.LB_Self()
	wg := sync.WaitGroup{}
	wg.Add(len(self.UPFInfos))
	for _, upf := range self.UPFInfos {
		go func(upf *lb_context.UpfContext) {
			updateCost := true
			if self.LoadBalancerType == LoadBalancerType_DP {
				updateCost = false
			}
			lb_client.GetMonitorData(upf, self.Topology.Delta, updateCost)
			wg.Done()
		}(upf)
	}
	wg.Wait()
}

func UpdatePathRemainRate(pathlist []up_topology.Path) {
	if pathlist == nil {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(pathlist))
	for _, path := range pathlist {
		go func(path up_topology.Path) {
			path.UpdateRemainRate()
			wg.Done()
		}(path)
	}
	wg.Wait()
}

func UpdatePathAll(pathlist []up_topology.Path) {
	wg := sync.WaitGroup{}
	wg.Add(len(pathlist))
	for _, path := range pathlist {
		go func(path up_topology.Path) {
			defer wg.Done()
			path.UpdateRemainRate()
			threshold := lb_context.LB_Self().Topology.Granularity
			if *path.RemainRate < threshold {
				logger.UtilLog.Debugf("path[%d] overload", path.Id)
				*path.Overload = true
				*path.Cost = 10000.0 // 10s
				return
			}
			*path.Overload = false
			if DP_ShortestPathThreshlod == -2 {
				dpUpdateCost(path)
			} else {
				path.UpdateCost()
			}
		}(path)
	}
	wg.Wait()
}

func ModifyPathRemainRate(path up_topology.Path, op RemainRateOperation, value float64) {
	switch op {
	case AddRemainRate:
		for _, rate := range path.RemainRates {
			*rate += value
		}
	case MinusRemainRate:
		for _, rate := range path.RemainRates {
			*rate -= value
		}
	}
}
