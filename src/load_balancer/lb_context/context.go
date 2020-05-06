package lb_context

import (
	"dynamicpath/lib/MonitorInfo"
	"dynamicpath/lib/loadbalancer_api"
	"dynamicpath/lib/up_topology"
	"net/http"
	"sync"
	"time"
)

var lbContext = LoadBalancerContext{}

var Wg sync.WaitGroup
var Mtx sync.Mutex

func init() {
	lbContext.Graph = up_topology.NewGraph()
	lbContext.UPFInfos = make(map[string]*UpfContext)
}

type LoadBalancerContext struct {
	Graph *up_topology.Graph
	Host  string
	// SmfUri            string
	// SmfClient         *http.Client
	Topology          *Topology
	QueryMonitorTimer *time.Timer
	UPFInfos          map[string]*UpfContext //NodeId->UpfContext
	LoadBalancerType  int
}

type Topology struct {
	SrcIp           string
	DstIp           string
	Delta           float64
	Granularity     float64 // Mbps that every pduSession needed
	PathThresh      int
	PathListAll     []up_topology.Path
	PathListBest    []up_topology.Path
	PathListWorst   []up_topology.Path
	RemainRateWorst float64
	RemainRateBest  float64
	RefreshTimer    *time.Timer
	PathInfos       map[int]*PathInfos //PathId
	RoundRobinCnt   int
}

type UpfContext struct {
	HttpServerUri  string
	HttpClient     *http.Client
	TopologyNode   *up_topology.Node
	InterfaceInfos map[string]InterfaceInfo     // eth -> interfaceInfo for Load(remain bandwidth)
	RemoteIpMaps   map[string]string            //remoteIp -> eth Name
	RemoteIpToEdge map[string]*up_topology.Edge // for Cost(packetloss and delay)
}

type PathInfos struct {
	PathId         int
	Path           up_topology.Path
	PduSessionInfo map[string]*PduSessionContext //supi-pduSessionid -> context
}

type PduSessionContext struct {
	Supi        string
	SessionInfo loadbalancer_api.SessionInfo
}
type InterfaceInfo struct {
	TotalBandwidth MonitorInfo.BandwidthInfo
	EthName        string
}

func NewUpfContext(serverUri string, node *up_topology.Node) *UpfContext {
	return &UpfContext{
		HttpServerUri:  serverUri,
		TopologyNode:   node,
		RemoteIpMaps:   make(map[string]string),
		InterfaceInfos: make(map[string]InterfaceInfo),
	}
}
func GetIfaceRxName(name string) string {
	return name + "-RX"
}
func GetIfaceTxName(name string) string {
	return name + "-TX"
}

func LB_Self() *LoadBalancerContext {
	return &lbContext
}
