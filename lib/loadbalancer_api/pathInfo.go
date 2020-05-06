package loadbalancer_api

type PathInfo struct {
	PathId    int        `json:"pathId"`
	EdgeInfos []EdgeInfo `json:"edgeInfos,omitempty"`
}
type EdgeInfo struct {
	StartIp     string `json:"startIp"`
	StartNodeId string `json:"startNodeId"`
	EndIp       string `json:"endIp"`
	EndNodeId   string `json:"endNodeId"`
}

type WorstPathList struct {
	PathList []int `json:"pathList"`
}
