package factory

type UpfInfo struct {
	NodeID string `yaml:"nodeId"`

	HttpUri string `yaml:"httpUri"`

	// LinkGnbInfo *LinkGnbInfo `yaml:"linkGnbInfo,omitempty"`

	// UserPlaneIps []string `yaml:"upIps"`

	// EdgeInfos []EdgeInfo `yaml:"edgeInfos"`
}
