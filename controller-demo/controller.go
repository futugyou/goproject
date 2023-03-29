package main

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	indexer  cache.Indexer                   // Indexer 的引⽤
	queue    workqueue.RateLimitingInterface //Workqueue 的引⽤
	informer cache.Controller                // Informer 的引⽤
}

// 将 Workqueue、Informer、Indexer 的引⽤作为参数返回⼀个新的 Controller
func NewController(queue workqueue.RateLimitingInterface, indexer cache.
	Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()
	klog.Info("Starting pod controller")
	// 启动 Informer 线程，Run 函数做两件事情 ：第⼀，运⾏⼀个 Reflector，并从 ListerWatcher
	// 中获取对象的通知放到队列中（Delta Queue）；第⼆，从队列中取出对象并处理该对象相关业务
	go c.informer.Run(stopCh)
	// 等待缓存同步队列
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("time out waitng for caches to sync"))
		return
	}
	// 启动多个 Worker 线程处理 Workqueue 中的 Object
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	// 启动⽆限循环，接收并处理消息
	for c.processNextItem() {
	}
}

// 从 Workqueue 中获取对象，并打印信息。
func (c *Controller) processNextItem() bool {
	key, shutdown := c.queue.Get()
	// 退出
	if shutdown {
		return false
	}
	// 标记此 Key 已经处理
	defer c.queue.Done(key)
	// 打印 Key 对应的 Object 的信息
	err := c.syncToStdout(key.(string))
	c.handleError(err, key)
	return true
}

func (c *Controller) handleError(err error, key interface{}) {
}

// 获取 Key 对应的 Object，并打印相关信息
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)

	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		fmt.Printf("Pod %s does not exist anymore\n", key)
	} else {
		fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*corev1.Pod).GetName())
	}

	return nil
}
