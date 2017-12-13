package kubeResourceProbe

import (
	"github.com/ericchiang/k8s"
	"github.com/golang/glog"
	"context"
	"os"
	"time"
	"flag"
	"sync"
	"github.com/ericchiang/k8s/watch/versioned"
	apiv1 "github.com/ericchiang/k8s/api/v1"
)

type  (
	ResourceProbe struct{
		client *k8s.Client
		once sync.Once
	}
    WatchableResources struct{
    	Configmaps []string
    	Secrets []string
    	ConfigmapChangeHandler handler
		SecretChangeHandler handler
		NS string
	}
	ResourceWatcher struct {
		confWatcher *ConfWatcher
		SecrtWatcher *SecrtWatcher
	}
	EvtWatcher interface {
        Close() error
	}
	ConfWatcher interface {
		EvtWatcher
		Next() (*versioned.Event, *apiv1.ConfigMap, error)
	}
	SecrtWatcher interface {
		EvtWatcher
		Next() (*versioned.Event, *apiv1.Secret, error)
	}
	handler func([]map[string]string)
)

const (
	CONF = iota
	SECRT
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

func (pp *ResourceProbe) WatchResource(resources *WatchableResources) {
	flag.Parse()
	pp.once.Do(func(){
		pp.client = pp.initClient()
	})
	go pp.startWatch(resources)
}

func (pp *ResourceProbe) startWatch(resources *WatchableResources){
	defer func() {
		if err := recover(); err != nil {
			glog.Error("Error occued in watch resource:", err)
			pp.wait()
			//rewatch
			pp.startWatch(resources)
		}
	}()
reWatch:
	confWatcher, err1 := pp.watchConfigmaps(resources.NS)
	secrtWatcher, err2 := pp.watchSecrets(resources.NS)
	if err1 != nil || err2 != nil {
		confWatcher.Close()
		secrtWatcher.Close()
		pp.wait()
		goto reWatch
	} else {
        evtWatcher1 := ConfWatcher(confWatcher)
		evtWatcher2 := SecrtWatcher(secrtWatcher)
		watchers := ResourceWatcher{
			&evtWatcher1,
			&evtWatcher2,
		}
		pp.watchResource(resources, watchers)
	}
}

func (pp *ResourceProbe) watchResource(resources *WatchableResources, watcher ResourceWatcher) {
	confWatcher := (*watcher.confWatcher)
	SecrtWatcher := (*watcher.SecrtWatcher)
	defer confWatcher.Close()
	defer (*watcher.SecrtWatcher).Close()

infiniteWar:
	evt1, got1, err1 := (*watcher.confWatcher).Next();
	evt2, got2, err2 := (*watcher.SecrtWatcher).Next();

	if err1 == nil {
		pp.processEvt(*evt1.Type, *got1.Metadata.Name, got1.Data, resources.Configmaps, resources.ConfigmapChangeHandler)
	}
	if err2 == nil {
		//slice performance downgrade, acceptable
		data := make(map[string]string)
		for k, v := range got2.Data {
			data[k] = string(v)
		}
		pp.processEvt(*evt2.Type, *got2.Metadata.Name, data, resources.Secrets, resources.SecretChangeHandler)
	}
    if err1 != nil {
		glog.Error("Error occured:", err1)
		confWatcher.Close()
		confWatcher, _ = pp.watchConfigmaps(resources.NS)
	}
	if err2 != nil {
		glog.Error("Error occured:", err2)
		SecrtWatcher.Close()
		SecrtWatcher, _ = pp.watchSecrets(resources.NS)
	}
	pp.wait()
	goto infiniteWar
}

func (pp *ResourceProbe) processEvt(eventType string, name string, data map[string]string, resources []string, handler handler) {
	if eventType == k8s.EventModified {
		gotFromChange := make([]map[string]string, len(resources))
		for _, cm := range resources {
			if name == cm {
				gotFromChange = append(gotFromChange, data)
				glog.Infof("resource %s update captured!", cm)
			}
		}
		if len(gotFromChange) > 0 {
			handler(gotFromChange)
		}
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

func (pp *ResourceProbe) watchSecrets(ns string) (*k8s.CoreV1SecretWatcher, error){
	CoreV1SecretWatcher, err := pp.client.CoreV1().WatchSecrets(context.Background(), ns)
	if err != nil {
		glog.Error("Watch secret Failed:", err)
	} else {
		glog.Error("Succeed in watching secrets...")
	}
	return CoreV1SecretWatcher, err
}

func (*ResourceProbe) wait() {
	time.Sleep(1 * time.Second)
}

