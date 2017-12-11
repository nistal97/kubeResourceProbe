package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
	"io/ioutil"
	"fmt"
	"github.com/ghodss/yaml"
)

type ResourceProbe struct{
	confFilePath string
	client *k8s.Client
}

func (pp *ResourceProbe) initClient() (*k8s.Client){
	_client, err := k8s.NewInClusterClient()
	if err != nil {
		glog.Error(err)
		panic("Failed to get client!")
	}
	return _client
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






