package node

import (
	"fmt"
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api/errors"
	"k8s.io/client-go/1.4/pkg/api/resource"
	"k8s.io/client-go/1.4/pkg/api/unversioned"
	api "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	client          *kubernetes.Clientset
	nodeIp          string
	updateFrequency int
	cores           int
	memory          int
	maxPods         int
}

type Node struct {
	Config
	pods map[string]api.Pod
}

func generateNodeIp(ipStr string, no int) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	ipInt := 0
	for _, a := range strings.Split(ipStr, ".") {
		num, err := strconv.Atoi(a)
		if err != nil {
			return ""
		}
		ipInt = ipInt*256 + num
	}

	ipInt += no
	buf := make([]string, 4)
	for i := 0; i < 4; i++ {
		a := ipInt % 256
		buf[3-i] = strconv.Itoa(a)
		ipInt = ipInt / 256
	}
	return strings.Join(buf, ".")
}

func NewConfig(conf conf.Config, no int) (*Config, error) {
	hostIp := generateNodeIp(conf.Ip, no)
	if len(hostIp) == 0 {
		return nil, fmt.Errorf("generateNodeIp failed")
	}
	config, err := clientcmd.BuildConfigFromFlags(conf.Apiserver, "")
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	nodeConfig := Config{
		client: client,
		nodeIp: hostIp,
	}
	return &nodeConfig, nil
}

func (n *Node) setNodeCapcity(node *api.Node) error {
	node.Status.Capacity = api.ResourceList{
		api.ResourceCPU: *resource.NewMilliQuantity(
			int64(n.cores*1000),
			resource.DecimalSI),
		api.ResourceMemory: *resource.NewQuantity(
			int64(n.memory),
			resource.BinarySI),
		api.ResourcePods: *resource.NewQuantity(
			int64(n.maxPods),
			resource.DecimalSI),
	}
	return nil
}

func (n *Node) setNodeReadyCondition(node *api.Node) error {
	currentTime := unversioned.NewTime(time.Now())
	newNodeReadyCondition := api.NodeCondition{
		Type:              api.NodeReady,
		Status:            api.ConditionTrue,
		Reason:            "KubeletReady",
		Message:           "kubelet is posting ready status",
		LastHeartbeatTime: currentTime,
	}

	readyConditionUpdated := false
	for i, condition := range node.Status.Conditions {
		if condition.Type != api.NodeReady {
			continue
		}
		if condition.Status != api.ConditionTrue {
			newNodeReadyCondition.LastTransitionTime = currentTime
		} else {
			newNodeReadyCondition.LastTransitionTime = condition.LastTransitionTime
		}
		node.Status.Conditions[i] = newNodeReadyCondition
		readyConditionUpdated = true
		break
	}

	if !readyConditionUpdated {
		node.Status.Conditions = append(node.Status.Conditions, newNodeReadyCondition)
	}
	return nil
}

func (n *Node) setNodeStatus(node *api.Node) error {
	n.setNodeCapcity(node)
	n.setNodeReadyCondition(node)
	return nil
}

func (n *Node) registerToApiserver() bool {
	node := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: n.nodeIp,
			Labels: map[string]string{
				"kubernetes.io/hostname":  n.nodeIp,
				"beta.kubernetes.io/os":   runtime.GOOS,
				"beta.kubernetes.io/arch": runtime.GOARCH,
			},
		},
		Spec: api.NodeSpec{
			Unschedulable: false,
		},
	}

	n.setNodeStatus(node)

	succ := false
	for i := 0; i < 5; i++ {
		if _, err := n.client.Core().Nodes().Create(node); err != nil {
			if !errors.IsAlreadyExists(err) {
				fmt.Printf("create node failed,", err)
				time.Sleep(time.Second)
			} else {
				n.client.Core().Nodes().Delete(n.nodeIp, nil)
			}
			continue
		}
		succ = true
		break
	}

	return succ
}

func (n *Node) heartBeat() {
	for i := 0; i < 5; i++ {
		node, err := n.client.Core().Nodes().Get(n.nodeIp)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		n.setNodeStatus(node)
		if _, err := n.client.Core().Nodes().Update(node); err != nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
}

func (n *Node) syncNodeStatus() {
	n.registerToApiserver()

	for {
		n.heartBeat()
		time.Sleep(time.Duration(n.updateFrequency) * time.Second)
	}
}

func (n *Node) Run() {
	go n.syncNodeStatus()
}
