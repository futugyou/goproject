package main

import (
	"flag"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

func main() {
	var kubeconfig string
	var master string
	// 从外部获取集群信息 (kube.config)
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	// 获取集群 master 的 url
	flag.StringVar(&master, "master", "", "master url")
	// 读取构建 config
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	// 创建 k8s Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	// 从指定的客户端、资源、命名空间和字段选择器创建⼀个新的 List-Watch
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(),
		"pods", v1.NamespaceDefault, fields.Everything())
	// 构造⼀个具有速率限制排队功能的新的 Workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	// 创建 Indexer 和 Informer
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{},
		0, cache.ResourceEventHandlerFuncs{
			//当有Pod创建时，根据Delta Queue弹出的Object⽣成对应的Key，并加⼊Workqueue中。
			//此处可以根据 Object 的⼀些属性进⾏过滤
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					queue.Add(key)
				}
			},
			//Pod 删除操作
			DeleteFunc: func(obj interface{}) {
				// 在⽣成 Key 之前检查对象。因为资源删除后有可能会进⾏重建等操作，如果监听时错过
				// 了删除信息，会导致该条记录是陈旧的
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)

				if err == nil {
					queue.Add(key)
				}
			},
		}, cache.Indexers{})
	// 创建新的 Controller
	controller := NewController(queue, indexer, informer)
	stop := make(chan struct{})
	defer close(stop)
	// 启动 Controller
	go controller.Run(1, stop)
	select {}
}
