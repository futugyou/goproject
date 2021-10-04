package controller

func Run(stopCh <-chan struct{}) error {
	run := func(stopCh <-chan struct{}) {
		err := StartController(stopCh)
		if err != nil {
			glog.Fatalf("error running service controllers: %v", err)
		}
		select {}
	}
	///忽略leader选举的相关逻辑
	run(stopCh)

	panic("unreachable")
}

func StartController(stopCh <-chan struct{}) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		glog.Fatalf("error building kubernetes config:%s", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	factory := informers.NewSharedInformerFactory(kubeClient, 0)

	podInformer := factory.Core().V1().Pods()
	pc := controller.NewPodController(kubeClient, podInformer, "k8s-cluster")
	go pc.Run(stopCh)

	factory.Start(stopCh)
	return nil
}
