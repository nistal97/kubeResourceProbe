package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
)

type ResourceProbe struct{
}

func (*ResourceProbe) initClient() (*k8s.Client){
	client, err := k8s.NewInClusterClient()
	if err != nil {
		glog.Error("Failed to init client!")
	} else {
		glog.Info("Succeed initing client!")
	}
	return client
}

func (pp *ResourceProbe) ListNodes() {
	nodes, _ := pp.initClient().CoreV1().ListNodes(context.Background())
	for _, node := range nodes.Items {
		//fmt.Println(node.String())
		glog.Info("name=%q schedulable=%t memory:%s/%s\n", *node.Metadata.Name, !*node.Spec.Unschedulable,
			*node.Status.Allocatable["memory"].String_, *node.Status.Capacity["memory"].String_)
	}
}

func (*ResourceProbe) WatchResource() {

}






