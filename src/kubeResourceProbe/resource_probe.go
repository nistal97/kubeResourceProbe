package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
    "os"
)

type ResourceProbe struct{
	confFilePath string
	client *k8s.Client
}

func Init() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default")
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
	pods, _ := pp.InitClient().CoreV1().ListPods(context.Background(), "app-ns")
	glog.Info(len(pods.Items))
}

func (*ResourceProbe) WatchResource() {

}






