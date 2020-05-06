package lb_util

import (
	"dynamicpath/lib/up_topology"
	"dynamicpath/src/load_balancer/factory"
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/logger"
)

func InitTopology() {
	config := factory.LBConfig
	self := lb_context.LB_Self()
	self.Host = config.Host
	self.Topology = &lb_context.Topology{
		SrcIp:       config.GnbIp,
		DstIp:       config.DnIp,
		Delta:       config.Delta,
		Granularity: config.Granularity,
		PathInfos:   make(map[int]*lb_context.PathInfos),
	}
	graph := self.Graph
	srcNode := graph.AddNode(config.GnbIp)
	dstNode := graph.AddNode(config.DnIp)
	for _, upfInfo := range config.UpfInfos {
		node := graph.AddNode(upfInfo.NodeID)
		upf := lb_context.NewUpfContext(upfInfo.HttpUri, node)
		self.UPFInfos[upfInfo.NodeID] = upf
	}
	for _, edgeInfo := range config.EdgeInfos {
		aNode, aExist := graph.Nodes[edgeInfo.SideInfos[0].NodeID]
		bNode, bExist := graph.Nodes[edgeInfo.SideInfos[1].NodeID]
		if !aExist || !bExist {
			logger.UtilLog.Panicf("NodeID %s or %s is not exist in Graph", edgeInfo.SideInfos[0].NodeID, edgeInfo.SideInfos[1].NodeID)
			return
		}
		aIp := edgeInfo.SideInfos[0].Ip
		bIp := edgeInfo.SideInfos[1].Ip
		aEdge, bEdge := graph.AddEdge(aNode, bNode, aIp, bIp)
		if aNode.NodeId == srcNode.NodeId || aNode.NodeId == dstNode.NodeId {
			upf := self.UPFInfos[bNode.NodeId]
			if upf.RemoteIpToEdge == nil {
				upf.RemoteIpToEdge = make(map[string]*up_topology.Edge)
			}
			upf.RemoteIpToEdge[aIp] = bEdge
		} else {
			upf := self.UPFInfos[aNode.NodeId]
			if upf.RemoteIpToEdge == nil {
				upf.RemoteIpToEdge = make(map[string]*up_topology.Edge)
			}
			upf.RemoteIpToEdge[bIp] = aEdge
		}
	}
	self.Topology.PathListAll = graph.GetAcyclicAllPath(*srcNode, *dstNode)
	self.Topology.PathListAll = delRedundantPath(graph, config.GnbIp, config.DnIp, self.Topology.PathListAll)
	// PathMergeSort(self.Topology.PathListAll)
	// for i, path := range self.Topology.PathListAll {
	// 	logger.UtilLog.Debugf("%d: %v, Cost-%.2f", i, path.NodeIdList, *path.Cost)
	// }
	// srcNode.RemainRates = nil
	// dstNode.RemainRates = nil

}
