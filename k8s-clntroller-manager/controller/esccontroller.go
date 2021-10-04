package controller

import (
	"fmt"
	"reflect"
	"sync"
)

type EcsBindingController struct {
	kubeClient              *versiond.Clientset
	clusterName             string
	ecsbingdingLister       list_and_watch.EcsBindingLister
	ecsbingdingListerSynced cache.InformerSynced
	broadcaster             record.EventBroadcaster
	recorder                record.EventRecorder

	ecsQueue workqueue.DelayingInterface
	lock     sync.Mutex
}

func NewEcsBindingController(kubeclient *versiond.Clientset, informer list_and_watch.EcsBindingInformer, clusterName string) *EcsBindingController {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "ecsbinding_controller"})

	ec := &EcsBindingController{
		kubeClient:              kubeclient,
		clusterName:             clusterName,
		ecsbingdingLister:       informer.Lister(),
		ecsbingdingListerSynced: informer.Informer().HasSynced,
		broadcaster:             eventBroadcaster,
		recorder:                recorder,

		ecsQueue: workqueue.NewNamedDelayingQueue("EcsBinding"),
	}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ecsbindingg := obj.(*ecsV1.EcsBinding)
			fmt.Printf("controller: Add event, ecsbinding [%s]\n", ecsbindingg.Name)
			ec.syncEcsbinding(ecsbindingg)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ecs1, ok1 := oldObj.(*ecsV1.EcsBinding)
			ecs2, ok2 := newObj.(*ecsV1.EcsBinding)
			if ok1 && ok2 && !reflect.DeepEqual(ecs1, ecs2) {
				fmt.Printf("controller: Update event, ecsbinding [%s]\n", ecs1.Name)
				ec.syncEcsbinding(ecs1)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ecsbindingg := obj.(*ecsV1.EcsBinding)
			fmt.Printf("controller: Delete event, ecsbinding [%s]\n", ecsbindingg.Name)
			ec.syncEcsbinding(ecsbindingg)
		},
	})

	return ec
}
