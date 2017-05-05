package node

import (
	"k8s.io/client-go/1.5/pkg/api/v1"
	"github.com/golang/glog"
	"time"
	"k8s.io/client-go/1.5/pkg/types"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api/errors"
	"k8s.io/client-go/1.5/pkg/api"
)

const ActTimeOut  = 10

const PodDelete v1.PodPhase = "Delete"

type Act struct {
	pod *v1.Pod
	exceptedStatus v1.PodPhase
}

type PodStatus struct {
	pod *v1.Pod
	exceptedStatus v1.PodPhase
	nextActTime time.Time
}

type StatusManager struct {
	pods map[types.UID] *PodStatus
	updates chan PodUpdate
	act chan Act
	client   *kubernetes.Clientset
}

func NewPodStatusManager(client *kubernetes.Clientset, updates chan PodUpdate) *StatusManager {
	statusManager := &StatusManager{}
	statusManager.client = client
	statusManager.updates = updates
	statusManager.pods = make(map[types.UID] *PodStatus)
	statusManager.act = make(chan Act, 20)
	return statusManager
}

func (m *StatusManager) processPodUpdates(update PodUpdate)  {
	cur := update.cur
	switch update.uType {
	case Create:
		m.pods[update.cur.UID] = &PodStatus{cur, v1.PodRunning, time.Now().Add(time.Duration(ActTimeOut) * time.Second)}
		glog.V(4).Infof("add pod:%s to status manager", cur.Name)
		m.act <- Act{cur, v1.PodRunning}
	case Update:
		old := update.old
		if old.UID != cur.UID {
			glog.Fatalf("pod updated, uid changed?,old:%s,new:%s", old.UID, cur.UID)
		}
		status, ok := m.pods[update.cur.UID]
		if !ok {
			glog.Fatalf("rec pod updated event, but no pods in status manager,update:%+v ", update)
		}
		status.pod = cur
		if cur.DeletionTimestamp != nil {
			m.act <- Act{cur, PodDelete}
		}
	case Delete:
		delete(m.pods, cur.UID)
		glog.V(4).Infof("delete pod:%s from status manager", cur.Name)
	}
}

func (status *PodStatus) check(now time.Time) *Act {
	if status.nextActTime.After(now) {
		return nil
	}
	if status.pod.Status.Phase != status.exceptedStatus {
		return &Act{status.pod, status.exceptedStatus}
	} else if status.pod.DeletionTimestamp != nil {
		return &Act{status.pod, PodDelete}
	}
	return nil
}

func (m *StatusManager) check() {
	for _, status := range m.pods {
		now := time.Now()
		if act := status.check(now); act != nil {
			glog.V(4).Infof("pod:%s, send a action,%+v,%s,nextActtime:%v", act.pod.Name, act, status.pod.Status.Phase, status.nextActTime)
			status.nextActTime = now.Add(time.Duration(ActTimeOut) * time.Second)
			m.act <- *act
		}
	}
}

func (m *StatusManager) updatePodStatus(pod *v1.Pod, phase v1.PodPhase) error {
	if phase == PodDelete {
		glog.V(4).Infof("try to delete a pod:%s", pod.Name)
		return m.client.Core().Pods(pod.Namespace).Delete(pod.Name, api.NewDeleteOptions(0))
	}
	glog.V(4).Infof("try to update pod:%s status:%s", pod.Name, phase)

	real, err := m.client.Core().Pods(pod.Namespace).Get(pod.Name)
	if err != nil {
		return err
	}

	real.Status.Phase = phase
	_, err = m.client.Core().Pods(pod.Namespace).UpdateStatus(real)
	return err
}

func (m *StatusManager) updatePodStatuses() {
	for act := range m.act {
		for i := 0;i < 3; i++ {
			err := m.updatePodStatus(act.pod, act.exceptedStatus)
			if err != nil {
				glog.Infof("pod:%s update status failed,err:%v", act.pod.Name, err)
				if errors.IsNotFound(err) {
					// a litte track
					// found pod delete event not delived to simluator
					m.updates <- PodUpdate{Delete, act.pod, nil}
					break
				}
			} else {
				break
			}
		}
	}
}

func (m * StatusManager) Run() {

	for i := 0; i < 5; i++  {
		go m.updatePodStatuses()
	}

	t := time.NewTicker(time.Second)
	for {
		select {
		case up :=<-m.updates:
			m.processPodUpdates(up)
		case <-t.C:
			m.check()
		}
	}
}

