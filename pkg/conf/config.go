package conf

var (
	SimConfig Config
)

type Config struct {
	Apiserver   string
	NodeNum     int
	Ip          string
	NodeCores   int
	NodeMem     int
	NodeMaxPods int
	UpdateFrequency int
}
