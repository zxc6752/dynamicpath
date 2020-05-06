package loadbalancer_api

type InitPathRequest struct {
	PathInfos []PathInfo `json:"pathInfos"`
}

type PathListAll struct {
	PathInfos []PathInfo `json:"pathInfos"`
}
