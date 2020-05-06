package up_topology

// type NICInfo struct {
// 	Name       string
// RamainRate float64
// }

type Node struct {
	NodeId      string
	Edges       []*Edge
	RemainRates map[string]*float64 // eth Name -> NIC remain rate
}

type Edge struct {
	Cost        float64
	StartIp     string
	EndIp       string
	Start       *Node
	End         *Node
	ReverseEdge *Edge
}

type Path struct {
	Id          int
	Cost        *float64
	RemainRate  *float64
	Edges       []*Edge
	NodeIdList  []string
	RemainRates []*float64 //NIC remain rate
	Overload    *bool
}

type Graph struct {
	Nodes map[string]*Node // NodeId -> Node
	Edges []*Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

func (g *Graph) GetAcyclicAllPath(src, dst Node) (pathList []Path) {
	visited := make(map[string]*bool)
	for _, node := range g.Nodes {
		tmp := false
		visited[node.NodeId] = &tmp
	}
	edges := make([]*Edge, len(g.Nodes)-1)
	nodeIdList := make([]string, len(g.Nodes))
	getAcyclicAllPathUtil(src, dst, visited, nodeIdList, edges, 0, &pathList)
	return
}

func getAcyclicAllPathUtil(cur, dst Node, visited map[string]*bool, nodeIdList []string, edges []*Edge, edgeCnt int, pathList *[]Path) {
	*visited[cur.NodeId] = true
	nodeIdList[edgeCnt] = cur.NodeId
	// if Current Node == Destination then append Path to PathList
	if cur.NodeId == dst.NodeId {
		overload := false
		path := Path{
			Cost:       new(float64),
			RemainRate: new(float64),
			Edges:      append([]*Edge{}, edges[:edgeCnt]...),
			NodeIdList: append([]string{}, nodeIdList[:edgeCnt+1]...),
			Overload:   &overload,
		}
		*pathList = append(*pathList, path)
	} else {
		for _, edge := range cur.Edges {
			if !*visited[edge.End.NodeId] {
				edges[edgeCnt] = edge
				getAcyclicAllPathUtil(*edge.End, dst, visited, nodeIdList, edges, edgeCnt+1, pathList)
			}

		}
	}
	*visited[cur.NodeId] = false
}

func (g *Graph) AddNode(nodeId string) *Node {
	g.Nodes[nodeId] = &Node{
		NodeId:      nodeId,
		RemainRates: make(map[string]*float64),
	}
	return g.Nodes[nodeId]
}

func (g *Graph) AddEdge(a, b *Node, aIp, bIp string) (*Edge, *Edge) {
	edge1 := &Edge{
		Cost:    1,
		StartIp: aIp,
		EndIp:   bIp,
		Start:   a,
		End:     b,
	}
	edge2 := &Edge{
		Cost:    1,
		StartIp: bIp,
		EndIp:   aIp,
		Start:   b,
		End:     a,
	}
	edge1.ReverseEdge = edge2
	edge2.ReverseEdge = edge1
	a.Edges = append(a.Edges, edge1)
	// remainA, remainB := 0.0, 0.0
	// a.RemainRates[aIp] = &remainA
	// b.RemainRates[bIp] = &remainB
	b.Edges = append(b.Edges, edge2)
	g.Edges = append(g.Edges, edge1, edge2)
	return edge1, edge2
}

func (p Path) UpdateRemainRate() {
	*p.RemainRate = *p.RemainRates[0]
	for _, rate := range p.RemainRates[1:] {
		if *rate < *p.RemainRate {
			*p.RemainRate = *rate
		}
	}
}

func (p Path) UpdateCost() {
	*p.Cost = p.Edges[0].Cost
	for _, edge := range p.Edges[1:] {
		*p.Cost += edge.Cost
	}
}

// func (p Path) GetNodeIds() []string {
// 	if p.Edges == nil {
// 		return nil
// 	}
// 	ipList := []string{p.Edges[0].Start.NodeId}
// 	for _, edge := range p.Edges {
// 		ipList = append(ipList, edge.End.NodeId)
// 	}
// 	return ipList
// }

func (n *Node) SameNICEdgeExist(localIp string) bool {
	cnt := 0
	for _, edge := range n.Edges {
		if localIp == edge.StartIp {
			cnt++
		}
	}
	if cnt > 1 {
		return true
	}
	return false
}

func (e *Edge) UpdateCost(cost float64) {
	e.Cost = cost
	e.ReverseEdge.Cost = cost
}
