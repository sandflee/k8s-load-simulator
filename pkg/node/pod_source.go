package node

import (
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/tools/cache"
	"k8s.io/client-go/1.5/pkg/fields"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/util/wait"
	"github.com/golang/glog"
)

type PodUpdateType string
const (
	Create PodUpdateType = "Create"
	Update PodUpdateType = "Update"
	Delete PodUpdateType = "Delete"
)

type PodUpdate struct {
	uType PodUpdateType
	cur *v1.Pod
	old *v1.Pod
}

func NewSourceApiserver(c *kubernetes.Clientset, nodeName string, updates chan<- PodUpdate) {
	lw := cache.NewListWatchFromClient(c.CoreClient, "pods", api.NamespaceAll, fields.OneTermEqualSelector(api.PodHostField, nodeName))
	newSourceApiserverFromLW(lw, updates)
}

// newSourceApiserverFromLW holds creates a config source that watches and pulls from the apiserver.
func newSourceApiserverFromLW(lw cache.ListerWatcher, updates chan<- PodUpdate) {
	_, controller := cache.NewInformer(lw, &v1.Pod{}, 0,
	  	cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			updates <- PodUpdate{Create, pod, nil}
		},
		UpdateFunc: func(old, cur interface{}) {
			pod := cur.(*v1.Pod)
			old2 := old.(*v1.Pod)
			updates <- PodUpdate{Update, pod, old2}
		},
		DeleteFunc: func(obj interface{}) {
			if pod, ok := obj.(*v1.Pod); ok {
				updates <- PodUpdate{Delete, pod, nil}
			} else {
				glog.Warning("recv:%+v", obj)
			}
		},
	})
	controller.Run(wait.NeverStop)
}
