package main

import (
	"time"

	kubeinformers "k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	clientset "kensho.ai/kdeployment/pkg/generated/clientset/versioned"
	informers "kensho.ai/kdeployment/pkg/generated/informers/externalversions"
	kdInformers "kensho.ai/kdeployment/pkg/generated/informers/externalversions/distribution.kensho.ai/v1"
)

func main() {
	stopCh := make(chan struct{})
	klog.InitFlags(nil)

	kubeconfig := "/Users/zhangjinrui/.kube/config"

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Error("Cannot parse kubeconfig")
	}

	// defaultClient, err := kubernetes.NewForConfig(cfg)
	kdeploymentClient, err := clientset.NewForConfig(cfg)
	kdInformerFactory := informers.NewSharedInformerFactory(kdeploymentClient, time.Second*30)

	kubeClient, err := kubernetes.NewForConfig(cfg)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

	kdInformerMap := make(map[string]kdInformers.KDeploymentInformer)
	kdInformerMap["cluster1"] = kdInformerFactory.Distribution().V1().KDeployments()

	kdClientMap := make(map[string]clientset.Interface)
	kdClientMap["cluster1"] = kdeploymentClient

	deploymentInformerMap := make(map[string]appsinformers.DeploymentInformer)
	deploymentInformerMap["cluster1"] = kubeInformerFactory.Apps().V1().Deployments()

	kubeClientMap := make(map[string]kubernetes.Interface)
	kubeClientMap["cluster1"] = kubeClient

	kdController := NewKDController(kdInformerMap, kdClientMap, deploymentInformerMap, kubeClientMap)
	kdInformerFactory.Distribution().V1().KDeployments().Informer()

	kdInformerFactory.Start(stopCh)
	kubeInformerFactory.Start(stopCh)

	if err = kdController.Run(1, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
