package main

import (
	"time"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	clientset "kensho.ai/kdeployment/pkg/generated/clientset/versioned"
	informers "kensho.ai/kdeployment/pkg/generated/informers/externalversions"
)

func main() {
	stop := make(chan struct{})
	klog.InitFlags(nil)

	kubeconfig := "/Users/zhangjinrui/.kube/config"

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Error("Cannot parse kubeconfig")
	}

	// defaultClient, err := kubernetes.NewForConfig(cfg)
	kdeploymentClient, err := clientset.NewForConfig(cfg)

	kdInformerFactory := informers.NewSharedInformerFactory(kdeploymentClient, time.Second*30)
	kdInformerFactory.Distribution().V1().KDeployments().Informer()
	kdInformerFactory.Start(stop)
}
