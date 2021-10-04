package controller

import (
	"fmt"
	"time"
)

type NodeController struct {
	kubeClient       *kubernetes.Clientset  //用于给符合条件的node打标记用
	nodeLister       corelisters.NodeLister //用于获取被监控的node资源
	nodeListerSynced cache.InformerSynced
	nodesQueue       workqueue.DelayingInterface  //一个延时队列，用于记录需要controller的node的key
	cloudProvider    cloudproviders.CloudProvider //用于判定node是否符合条件打标记，此成员并非controller关键结构的成员
}

func NewNodeController(kubeClient *kubernetes.Clientset, nodeInformer coreinformers.NodeInformer, cp cloudproviders.CloudProvider) *NodeController {

	n := &NodeController{
		kubeClient:       kubeClient,
		nodeLister:       nodeInformer.Lister(),
		nodeListerSynced: nodeInformer.Informer().HasSynced,
		nodesQueue:       workqueue.NewNamedDelayingQueue("nodes"),
		cloudProvider:    cp,
	}

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(cur interface{}) {
			node := cur.(*v1.Node)
			fmt.Printf("controller: Add event, nodes [%s]\n", node.Name)
			n.syncNodes(node)
		},
		DeleteFunc: func(cur interface{}) {
			node := cur.(*v1.Node)
			fmt.Printf("controller: Delete event, nodes [%s]\n", node.Name)
			n.syncNodes(node)
		},
	})

	return n
}

func (n *NodeController) syncNodes(node *v1.Node) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(node)
	if err != nil {
		fmt.Printf("Couldn't get key for object %#v: %v \n", node, err)
		return
	}
	fmt.Printf("working queue add node %s\n", key)
	n.nodesQueue.Add(key)
}

func (n *NodeController) Run(stopCh <-chan struct{}) {

	defer runtime.HandleCrash()
	defer n.nodesQueue.ShutDown()

	fmt.Println("Starting service controller")
	defer fmt.Println("Shutting down service controller")

	if !cache.WaitForCacheSync(stopCh, n.nodeListerSynced) {
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
		return
	}

	for i := 0; i < WorkerCount; i++ {
		go wait.Until(n.worker, time.Second, stopCh)
	}

	<-stopCh
}

func (n *NodeController) worker() {
	for {
		func() {
			key, quit := n.nodesQueue.Get()
			if quit {
				return
			}
			defer n.nodesQueue.Done(key)
			err := n.handleNode(key.(string))
			if err != nil {
				fmt.Printf("controller: error syncing node, %v \n", err)
			}
		}()
	}
}

func (n *NodeController) handleNode(key string) error {
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	node, err := n.nodeLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		fmt.Printf("Node has been deleted %v \n", key)
		return nil
	case err != nil:
		fmt.Printf("Unable to retrieve node %v from store: %v \n", key, err)
		n.nodesQueue.AddAfter(key, time.Second*30)
		return err
	default:
		err := n.processNodeAddIntoCluster(node)
		if err != nil {
			n.nodesQueue.AddAfter(key, time.Second*30)
		}
		return err
	}
	return nil
}
