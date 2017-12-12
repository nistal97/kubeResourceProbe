package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
	"os"
	"time"
	"flag"
)

type  (
	ResourceProbe struct{
		client *k8s.Client
	}
    WatchableResources struct{
    	Configmaps []string
    	Secrets []string
    	ConfigmapChangeHandler func([]map[string]string)
		SecretChangeHandler func()
		NS string
	}
)

func Init() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	flag.Parse()
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

func (pp *ResourceProbe) WatchResource(resources *WatchableResources) {
	if pp.client == nil {
		pp.client = pp.initClient()
	}
	go pp.watchResource(resources)
}

func (pp *ResourceProbe) watchResource(resources *WatchableResources){
	defer func() {
		if err := recover(); err != nil {
			glog.Error("Error occued in watch resource:", err)
		}
	}()

	CoreV1ConfigMapWatcher, err := pp.watchConfigmaps(resources.NS)
	defer CoreV1ConfigMapWatcher.Close()
	if err != nil {
		glog.Error("Failed to watch configmaps, keep trying:", err)
	} else {
infiniteWar:
		if event, got, err := CoreV1ConfigMapWatcher.Next(); err != nil {
			glog.Error("Failed to get next watch event")
		} else {
			if *event.Type == k8s.EventModified {
				confs := make([]map[string]string, len(resources.Configmaps))
				for _, cm := range resources.Configmaps {
					if got.Metadata.GetName() == cm {
						confs = append(confs, got.Data)
						glog.Info("configmap %q update captured!", cm)
					}
				}
				if len(confs) > 0 {
					resources.ConfigmapChangeHandler(confs)
				}
			}
		}
		time.Sleep(1 * time.Second)
		goto infiniteWar
	}
}

func (pp *ResourceProbe) watchConfigmaps(ns string) (*k8s.CoreV1ConfigMapWatcher, error){
	CoreV1ConfigMapWatcher, err := pp.client.CoreV1().WatchConfigMaps(context.Background(), ns)
	if err != nil {
		glog.Error("Watch Configmaps Failed:", err)
	} else {
		glog.Error("Succeed in watching configmaps...")
	}
	return CoreV1ConfigMapWatcher, err
}




