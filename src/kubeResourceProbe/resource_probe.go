package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
    "os"
	"time"
)

type  (
	ResourceProbe struct{
		client *k8s.Client
	}
    WatchableResources struct{
    	configmaps []string
    	secrets []string
	}
)

func Init() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
}

func (*ResourceProbe) initClient() (*k8s.Client){
	client, err := k8s.NewInClusterClient()
	if err != nil {
		glog.Error("Failed to init client!", err)
	} else {
		glog.Info("Succeed initing client!")
	}
	return client
}

func (pp *ResourceProbe) WatchResource(resources WatchableResources) {
	if pp.client == nil {
		pp.client = pp.initClient()
	}
	CoreV1ConfigMapWatcher, err := pp.client.CoreV1().WatchConfigMaps(context.Background(), "app-ns")
	defer CoreV1ConfigMapWatcher.Close()
	if err != nil {
		glog.Error("Watch Configmaps Failed:", err)
	} else {
		glog.Info("Now watching configmaps...")
		go pp.watch(CoreV1ConfigMapWatcher)
	}
}

func (*ResourceProbe) watch(CoreV1ConfigMapWatcher *k8s.CoreV1ConfigMapWatcher){
	defer func() {
		if err := recover(); err != nil {
			glog.Error("Error occued in watch resource:", err)
		}
	}()
infiniteWatch:
	if event, _, err := CoreV1ConfigMapWatcher.Next(); err != nil {
		glog.Error("Failed to watch configmaps:", err)
	} else {
		if *event.Type == k8s.EventModified {
			glog.Info("configMap is modified..")
		}
	}
	time.Sleep(3 * time.Second)
	goto infiniteWatch
}





