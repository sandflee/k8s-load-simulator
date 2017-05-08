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
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"github.com/sandflee/k8s-load-simulator/pkg/node"
	"time"
	"github.com/golang/glog"
)

func Run() {
	glog.Info("run k8s simulator")
	for i := 0; i < conf.SimConfig.NodeNum; i++ {
		c, err := node.NewConfig(conf.SimConfig, i)
		if err != nil {
			glog.Fatal("create config failed,%v\n", err)
		}
		node := node.Node{
			Config: *c,
		}
		go node.Run()
	}

	cacher, err := node.NewNodeCacher(conf.SimConfig.Apiserver)
	if err != nil {
		glog.Fatal("create node cacher failed", err)
	}
	go cacher.Run()

	timer := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-timer.C:
			glog.Info("ping")
		}
	}
}
