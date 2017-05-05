package node

import (
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/tools/cache"
	"k8s.io/client-go/1.5/pkg/fields"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"github.com/golang/glog"
)

func NewSourceApiserver(c *kubernetes.Clientset, nodeName string, updates chan<- interface{}) {
	lw := cache.NewListWatchFromClient(c.CoreClient, "pods", api.NamespaceAll, fields.OneTermEqualSelector(api.PodHostField, nodeName))
	newSourceApiserverFromLW(lw, updates)
}

// newSourceApiserverFromLW holds creates a config source that watches and pulls from the apiserver.
func newSourceApiserverFromLW(lw cache.ListerWatcher, updates chan<- interface{}) {

	send := func(objs []interface{}) {
		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		glog.Info("recevice pods:%+v", pods)
		//updates <- kubetypes.PodUpdate{Pods: pods, Op: kubetypes.SET, Source: kubetypes.ApiserverSource}
	}
	cache.NewReflector(lw, &api.Pod{}, cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc), 0).Run()
}
