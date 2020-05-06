package factory

type Config struct {
	Period int `yaml:"period"`

	Delta float64 `yaml:"delta"`

	Granularity float64 `yaml:"granularity"`

	Host string `yaml:"host"`

	GnbIp string `yaml:"gnbIp"`

	DnIp string `yaml:"dnIp"`

	UpfInfos []UpfInfo `yaml:"upfInfos"`

	EdgeInfos []EdgeInfo `yaml:"edgeInfos"`

	Logger Logger `yaml:"logger"`

	LoadBalancerType int `yaml:"loadBalancerType"`
}
