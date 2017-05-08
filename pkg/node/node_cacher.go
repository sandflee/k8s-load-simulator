package node

import (
	"k8s.io/client-go/1.5/kubernetes"
	api "k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/tools/clientcmd"
	"github.com/golang/glog"
	"k8s.io/client-go/1.5/tools/cache"
	"k8s.io/client-go/1.5/pkg/util/wait"
	"time"
)

type NodeInfo struct {
	lastUpdateTime time.Time
	//intervals []time.Duration
	interval time.Duration
}

func (n *NodeInfo) updateTime(node *api.Node) {
	for _, condition := range node.Status.Conditions {
		if condition.Type == api.NodeReady {
			if condition.Status != api.ConditionTrue {
				glog.Infof("node:%s, ready condition is not true,%+v", node.Name, condition)
			}
			if n.lastUpdateTime != condition.LastHeartbeatTime.Time {
				n.interval = time.Now().Sub(condition.LastHeartbeatTime.Time)
				if n.interval > time.Second {
					glog.V(0).Infof("node:%s updated, heartbeart time from %v to %v, time used:%v", node.Name, n.lastUpdateTime, condition.LastHeartbeatTime.Time, n.interval)
				} else {
					glog.V(3).Infof("node:%s updated, heartbeart time from %v to %v, time used:%v", node.Name, n.lastUpdateTime, condition.LastHeartbeatTime.Time, n.interval)
				}
				n.lastUpdateTime = condition.LastHeartbeatTime.Time
			}
			break
		}
	}
}

type NodeCacher struct {
	nodeStore cache.Store
	controller *cache.Controller
	nodes map[string]*NodeInfo
}

func (n *NodeCacher) updateNodeInfo(node *api.Node, isDelete bool) {
	if isDelete {
		glog.V(1).Infof("node:%s,delete from node cacher", node.Name)
		delete(n.nodes, node.Name)
		return
	}
	info, ok := n.nodes[node.Name]
	if !ok {
		glog.V(1).Infof("node:%s,add to node cacher", node.Name)
		info = &NodeInfo{}
		n.nodes[node.Name] = info
	}
	info.updateTime(node)
}

func NewNodeCacher(apiserver string) (*NodeCacher, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, "")
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cacher := &NodeCacher{}

	lw := cache.NewListWatchFromClient(client.CoreClient, "nodes", api.NamespaceAll, nil)
	store, controller := cache.NewInformer(lw, &api.Node{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node := obj.(*api.Node)
				cacher.updateNodeInfo(node, false)
			},
			UpdateFunc: func(old, cur interface{}) {
				node := cur.(*api.Node)
				cacher.updateNodeInfo(node, false)
			},
			DeleteFunc: func(obj interface{}) {
				if node, ok := obj.(*api.Node); ok {
					cacher.updateNodeInfo(node, true)
				} else {
					glog.Warningf("node delete handler, recv:%v",obj)
				}
			},
		},
	)

	cacher.nodeStore = store
	cacher.controller = controller
	cacher.nodes = make(map[string]*NodeInfo)
	return cacher, nil
}

func (n *NodeCacher) Run() {
	n.controller.Run(wait.NeverStop)
}
