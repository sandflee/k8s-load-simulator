// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simulator

import (
	"fmt"
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"github.com/sandflee/k8s-load-simulator/pkg/node"
	"github.com/spf13/cobra"
	"time"
)

func DoRun(cmd *cobra.Command, args []string) {
	//defer glog.Flush()

	for i := 0; i < conf.SimConfig.NodeNum; i++ {
		c, err := node.NewConfig(conf.SimConfig, i)
		if err != nil {
			fmt.Printf("create config failed,%v\n", err)
			return
		}
		node := node.Node{
			Conf: *c,
		}
		go node.Run()
	}

	timer := time.NewTicker(time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Printf("timer comes")
		}
	}
}
