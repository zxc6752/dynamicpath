package MonitorInfo

type InterfacesInfo struct {
	BandwidthMaps map[string]BandwidthInfo `json:"bandwidthMaps"`
	RemoteIpMaps  map[string]string        `json:"remoteIpMaps,omitempty"`
}
type BandwidthInfo struct {
	Rx float64 `json:"rx"`
	Tx float64 `json:"tx"`
}

type ConnectionInfo struct {
	PacketLoss float64 `json:"loss"`  // 0~100%
	DelayTime  float64 `json:"delay"` // ms
}
type UPFMonitorData struct {
	// CpuUsage        float64                    `json:"cpuUsage"`    // 0~100%
	PacketRates     map[string]BandwidthInfo   `json:"packetRate"`  // interfaceName -> Mbps
	ConnectionInfos map[string]*ConnectionInfo `json:"connections"` // Link's Remote IP -> Traffic info
}
