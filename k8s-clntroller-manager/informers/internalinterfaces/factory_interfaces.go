package internalinterfaces

import "time"

type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	InformerFor(obj runtime.Object, newFunc NewInformerFunc) cache.SharedIndexInformer
}

type NewInformerFunc func(versiond.Interface, time.Duration) cache.SharedIndexInformer
