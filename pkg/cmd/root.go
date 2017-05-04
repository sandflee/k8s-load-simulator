package cmd

import (
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"github.com/sandflee/k8s-load-simulator/pkg/simulator"
	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8s-load-simulator",
	Short: "A low-level load simulater tool for k8s",
	Run:   simulator.DoRun,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&conf.SimConfig.Apiserver, "apiserver", "http://127.0.0.1:8080", "apiserver address")
	RootCmd.PersistentFlags().IntVar(&conf.SimConfig.NodeNum, "nodeNum", 1, "Total number of mockNode")
	RootCmd.PersistentFlags().StringVar(&conf.SimConfig.Ip, "nodeIp", "127.0.0.1", "the first mock node ip address")
}
