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

package main

import (
	"flag"
	"github.com/sandflee/k8s-load-simulator/pkg/simulator"
	"github.com/golang/glog"
	_ "github.com/sandflee/k8s-load-simulator/pkg/cmd"
	"net/http"
	_ "net/http/pprof"
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"strconv"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	go simulator.Run()

	glog.Infof("config:%+v", conf.SimConfig)

	err := http.ListenAndServe("0.0.0.0:" + strconv.Itoa(conf.SimConfig.PprofPort), nil)
	if err != nil {
		glog.Fatal("http listen failed,", err)
	}
}
