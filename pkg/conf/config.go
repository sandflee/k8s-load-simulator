package conf

var (
	SimConfig Config
)

type Config struct {
	Apiserver string
	NodeNum   int
	Ip        string
}
