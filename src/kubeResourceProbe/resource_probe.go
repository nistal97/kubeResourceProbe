package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
	"fmt"
    "os"
)

type ResourceProbe struct{
	confFilePath string
	client *k8s.Client
}

func Init() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc.cluster.local")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

func (*ResourceProbe) InitClient() (*k8s.Client){
	client, err := k8s.NewInClusterClient()
	if err != nil {
		glog.Error("Failed to init client!", err)
	} else {
		glog.Info("Succeed initing client!")
	}

	return client
}

func (pp *ResourceProbe) ListNodes() {
	nodes, _ := pp.InitClient().CoreV1().ListNodes(context.Background())
	for _, node := range nodes.Items {
		//fmt.Println(node.String())
		glog.Info("name=%q schedulable=%t memory:%s/%s\n", *node.Metadata.Name, !*node.Spec.Unschedulable,
			*node.Status.Allocatable["memory"].String_, *node.Status.Capacity["memory"].String_)
	}
}

func (*ResourceProbe) WatchResource() {

}






