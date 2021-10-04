package informers

import (
	"reflect"
	"time"
)

type SharedInformerFactory interface {
	internalinterfaces.SharedInformerFactory
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool

	EcsBind() ecsbind.Interface
}
type EcsBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() demov1.EcsBindingLister
}

type ecsBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func (e *ecsBindingInformer) Informer() cache.SharedIndexInformer {
	return e.factory.InformerFor(&ecsv1.EcsBinding{}, e.defaultInformer)
}

func (e *ecsBindingInformer) defaultInformer(client versiond.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredEcsBindingInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, e.tweakListOptions)
}

func (e *ecsBindingInformer) Lister() demov1.EcsBindingLister {
	return demov1.NewEcsBindingLister(e.Informer().GetIndexer())
}

func NewFilteredEcsBindingInformer(clientset *versiond.Clientset, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (object runtime.Object, e error) {
			if tweakListOptions != nil {
				tweakListOptions(&options)
			}
			return clientset.EcsV1().EcsBinding().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (i watch.Interface, e error) {
			if tweakListOptions != nil {
				tweakListOptions(&options)
			}
			return clientset.EcsV1().EcsBinding().Watch(options)
		},
	}, &ecsv1.EcsBinding{}, resyncPeriod, indexers)
}
