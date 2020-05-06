package MonitorInfo

type MonitorStartRequest struct {
	LinkIPInfo []MonitorIPInfo `json:"linkInfo"`
}

type MonitorIPInfo struct {
	IP     string `json:"ip"`
	Update bool   `json:"updated"`
}
