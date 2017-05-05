package cmd

import (
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"flag"
)


func init() {
	flag.StringVar(&conf.SimConfig.Apiserver, "apiserver", "http://127.0.0.1:8080", "apiserver address")
	flag.IntVar(&conf.SimConfig.NodeNum, "nodeNum", 1, "Total number of mockNode")
	flag.IntVar(&conf.SimConfig.NodeCores, "nodeCores", 16, "cpu capacity for node")
	flag.IntVar(&conf.SimConfig.NodeMem, "nodeMem", 32*1024, "mem capacity for node")
	flag.IntVar(&conf.SimConfig.NodeMaxPods, "nodeMaxPods", 100, "max pods that could running on nodes")
	flag.IntVar(&conf.SimConfig.UpdateFrequency, "heartbeat-interval", 10, "heartbeat-interval of seconds")
	flag.IntVar(&conf.SimConfig.PprofPort, "pprof-port", 6666, "pprof listen port")
	flag.StringVar(&conf.SimConfig.Ip, "nodeIp", "127.0.0.1", "the first mock node ip address")
}
