package lb_util

import (
	"dynamicpath/lib/up_topology"
	"fmt"
)

func delRedundantPath(graph *up_topology.Graph, srcIp, dstIp string, orignial []up_topology.Path) []up_topology.Path {
	sameNICEdges := make(map[string]string) // (StartNodeId-startIp)->(EndNodeId-EndIp)
	for _, edge := range graph.Edges {

		if edge.Start.NodeId == srcIp {
			continue
		}
		if edge.End.NodeId == dstIp {
			continue
		}
		if edge.Start.SameNICEdgeExist(edge.StartIp) && edge.End.SameNICEdgeExist(edge.EndIp) {
			key := fmt.Sprintf("%s-%s", edge.Start.NodeId, edge.StartIp)
			sameNICEdges[key] = fmt.Sprintf("%s-%s", edge.End.NodeId, edge.EndIp)
		}
	}
	if len(sameNICEdges) < 1 {
		for _, path := range orignial {
			path.UpdateCost()
		}
		return orignial
	}
	deletePathId := []int{}

	for i, path := range orignial {
		trace := make(map[string]bool)
		for _, edge := range path.Edges[1 : len(path.Edges)-1] {
			end := fmt.Sprintf("%s-%s", edge.End.NodeId, edge.EndIp)
			if trace[end] {
				deletePathId = append(deletePathId, i)
				break
			}
			start := fmt.Sprintf("%s-%s", edge.Start.NodeId, edge.StartIp)
			if endInfo, exist := sameNICEdges[start]; exist {
				trace[endInfo] = true
			}

		}
		path.UpdateCost()
		// fmt.Printf("%d: %v, Cost-%.2f\n", i, path.NodeIdList, *path.Cost)
	}

	newPathList := []up_topology.Path{}
	cnt := 0
	for _, index := range deletePathId {
		newPathList = append(newPathList, orignial[cnt:index]...)
		cnt = index + 1
	}
	newPathList = append(newPathList, orignial[cnt:]...)
	return newPathList
}
