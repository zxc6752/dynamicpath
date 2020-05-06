package loadbalancer_api

type SessionInfo struct {
	PduSessionId int `json:"pduSessionId"`
	PathID       int `json:"pathID"`
}
