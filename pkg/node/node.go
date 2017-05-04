package node

import (
	api "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"github.com/sandflee/k8s-load-simulator/pkg/conf"
	"net"
	"strings"
	"strconv"
	"fmt"
	"runtime"
	"k8s.io/client-go/1.4/pkg/api/errors"
	"time"
)

type Config  struct {
	client *kubernetes.Clientset
	nodeIp string
	update int
}

type Node struct {
	Conf Config
	pods map[string] api.Pod
}

func generateNodeIp(ipStr string, no int) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return  ""
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
		a := ipInt%256
		buf[3-i] = strconv.Itoa(a)
		ipInt = ipInt/256
	}
	return strings.Join(buf, ".")
}

func NewConfig(conf conf.Config, no int) (*Config, error) {
	hostIp := generateNodeIp(conf.Ip, no)
	if len(hostIp) == 0 {
		return nil, fmt.Errorf("generateNodeIp failed")
	}
	config, err := clientcmd.BuildConfigFromFlags(conf.Apiserver,"");
	if err != nil {
		return nil, err
	}
	client ,err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	nodeConfig := Config {
		client : client,
		nodeIp : hostIp,
	}
	return &nodeConfig, nil
}

func (n *Node) setNodeStatus(node *api.Node) *api.Node {
	return node
}

func (n *Node) registerToApiserver() bool {
	node := &api.Node{
		ObjectMeta: api.ObjectMeta{
			Name: n.Conf.nodeIp,
			Labels: map[string]string {
				"kubernetes.io/hostname": n.Conf.nodeIp,
				"beta.kubernetes.io/os": runtime.GOOS,
				"beta.kubernetes.io/arch": runtime.GOARCH,
			},
		},
		Spec: api.NodeSpec {
			Unschedulable: false,
		},
	}

	node = n.setNodeStatus(node)

	succ := false
	for i:=0; i < 5; i++ {
		if _, err := n.Conf.client.Core().Nodes().Create(node); err != nil {
			if !errors.IsAlreadyExists(err) {
				fmt.Printf("create node failed,", err)
				time.Sleep(time.Second)
				continue
			}
		}
		succ = true
		break
	}

	return succ
}

func (n *Node) heartBeat() {
	for i:=0; i < 5; i++ {
		node , err := n.Conf.client.Core().Nodes().Get(n.Conf.nodeIp)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		node = n.setNodeStatus(node)
		if _, err := n.Conf.client.Core().Nodes().Update(node); err != nil {
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
		time.Sleep(time.Duration(n.Conf.update) * time.Second)
	}
}

func (n *Node) Run() {
	go n.syncNodeStatus()
}
