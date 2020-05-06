package lb_util

import (
	"dynamicpath/lib/up_topology"
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/logger"
)

type PathListType string

const (
	WorstPath PathListType = "WorstPath"
	BestPath  PathListType = "BestPath"
)

func pathMerge(leftPaths, rightPaths []up_topology.Path) {
	leftLen, rightLen := len(leftPaths), len(rightPaths)
	leftIndex, rightIndex := 0, 0
	var sortedPath []up_topology.Path
	for i := 0; i < rightLen+leftLen; i++ {
		if leftIndex == leftLen {
			sortedPath = append(sortedPath, rightPaths[rightIndex:]...)
			break
		} else if rightIndex == rightLen {
			sortedPath = append(sortedPath, leftPaths[leftIndex:]...)
			break
		} else if *leftPaths[leftIndex].Cost < *rightPaths[rightIndex].Cost {
			sortedPath = append(sortedPath, leftPaths[leftIndex])
			leftIndex++
		} else {
			sortedPath = append(sortedPath, rightPaths[rightIndex])
			rightIndex++
		}
	}
	for i := 0; i < leftLen; i++ {
		leftPaths[i] = sortedPath[i]
	}
	for i := 0; i < rightLen; i++ {
		rightPaths[i] = sortedPath[leftLen+i]
	}
}

func dpUpdateCost(path up_topology.Path) {
	self := lb_context.LB_Self()
	*path.Cost = 0
	for _, edge := range path.Edges {
		upf := self.UPFInfos[edge.Start.NodeId]
		if upf != nil {
			ethName := upf.RemoteIpMaps[edge.EndIp]
			*path.Cost = *path.Cost + (1-*edge.Start.RemainRates[lb_context.GetIfaceTxName(ethName)]/upf.InterfaceInfos[ethName].TotalBandwidth.Tx)*100
		}
		upf = self.UPFInfos[edge.End.NodeId]
		if upf != nil {
			ethName := upf.RemoteIpMaps[edge.StartIp]
			*path.Cost = *path.Cost + (1-*edge.End.RemainRates[lb_context.GetIfaceRxName(ethName)]/upf.InterfaceInfos[ethName].TotalBandwidth.Rx)*100
		}
	}
}

func PathMergeSort(pathList []up_topology.Path) {
	if len(pathList) > 1 {
		middle := len(pathList) / 2
		leftPaths := pathList[:middle]
		rightPaths := pathList[middle:]
		PathMergeSort(leftPaths)
		PathMergeSort(rightPaths)
		pathMerge(leftPaths, rightPaths)
	}
}

func GetSubPathList(pathList []up_topology.Path, types PathListType) (subPath []up_topology.Path, mean float64) {
	if len(pathList) < 1 {
		return nil, 0.0
	}
	mean = PathRemainRateMean(pathList)
	switch types {
	case WorstPath:
		for _, path := range pathList {
			if *path.RemainRate <= mean {
				subPath = append(subPath, path)
			}
		}
	case BestPath:
		for _, path := range pathList {
			if *path.RemainRate >= mean {
				subPath = append(subPath, path)
			}
		}
	}
	return
}

func PathRemainRateMean(pathList []up_topology.Path) float64 {
	sum := 0.0
	// overLoadCnt := 0
	for _, path := range pathList {
		// if *path.Overload {
		// 	overLoadCnt++
		// 	continue
		// }
		sum = sum + *path.RemainRate
	}
	// n := len(pathList) - overLoadCnt
	// if n == 0 {
	// return 0.0
	// }
	return sum / float64(len(pathList))
}

func InitAllPath() {
	self := lb_context.LB_Self()
	for i, path := range self.Topology.PathListAll {
		for _, edge := range path.Edges {
			// Init Path Link To Node Remain Rate
			// Only Record consider Uplink for one path
			upf := self.UPFInfos[edge.Start.NodeId]
			if upf != nil {
				ethName := upf.RemoteIpMaps[edge.EndIp]
				path.RemainRates = append(path.RemainRates, edge.Start.RemainRates[lb_context.GetIfaceTxName(ethName)])
			}
			upf = self.UPFInfos[edge.End.NodeId]
			if upf != nil {
				ethName := upf.RemoteIpMaps[edge.StartIp]
				path.RemainRates = append(path.RemainRates, edge.End.RemainRates[lb_context.GetIfaceRxName(ethName)])
			}
		}
		path.Id = i
		self.Topology.PathInfos[i] = &lb_context.PathInfos{
			Path:           path,
			PduSessionInfo: make(map[string]*lb_context.PduSessionContext),
		}
		self.Topology.PathListAll[i] = path
		logger.UtilLog.Debugf("pathId: %d , nodes: %v", i, path.NodeIdList)
	}
}

func FindPathThreshold(pathList []up_topology.Path) int {
	middle := len(pathList)/2 - (1 - len(pathList)%2)
	if !*pathList[middle].Overload {
		return middle
	}
	left, right := 0, middle-1
	for left <= right {
		middle = (left + right) / 2
		if *pathList[middle].Overload {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}
	return right
}
