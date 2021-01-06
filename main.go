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

var clusterMap = map[string]string{
	"aks-0": "/Users/zhangjinrui/.kube/config-aks-0",
	"aks-1": "/Users/zhangjinrui/.kube/config-aks-1",
	"aks-2": "/Users/zhangjinrui/.kube/config-aks-2",
}

// var clusterMap = map[string]string{
// 	"aks-0": "/Users/zhangjinrui/.kube/config-backup",
// }

func main() {
	stopCh := make(chan struct{})
	klog.InitFlags(nil)

	kdInformerMap := make(map[string]kdInformers.KDeploymentInformer)
	kdClientMap := make(map[string]clientset.Interface)
	deploymentInformerMap := make(map[string]appsinformers.DeploymentInformer)
	kubeClientMap := make(map[string]kubernetes.Interface)

	var kdeploymentInformerList []informers.SharedInformerFactory
	var deploymentInformerList []kubeinformers.SharedInformerFactory

	var kdController *KDController

	for clusterName, configPath := range clusterMap {
		kubeconfig := configPath

		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			klog.Error("Cannot parse kubeconfig")
		}

		// defaultClient, err := kubernetes.NewForConfig(cfg)
		kdeploymentClient, err := clientset.NewForConfig(cfg)
		kdInformerFactory := informers.NewSharedInformerFactory(kdeploymentClient, time.Second*30)

		kubeClient, err := kubernetes.NewForConfig(cfg)
		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

		kdInformerMap[clusterName] = kdInformerFactory.Distribution().V1().KDeployments()
		kdClientMap[clusterName] = kdeploymentClient
		deploymentInformerMap[clusterName] = kubeInformerFactory.Apps().V1().Deployments()
		kubeClientMap[clusterName] = kubeClient

		kdeploymentInformerList = append(kdeploymentInformerList, kdInformerFactory)
		deploymentInformerList = append(deploymentInformerList, kubeInformerFactory)
	}

	kdController = NewKDController(kdInformerMap, kdClientMap, deploymentInformerMap, kubeClientMap)

	for _, informer := range kdeploymentInformerList {
		informer.Start(stopCh)
	}

	for _, informer := range deploymentInformerList {
		informer.Start(stopCh)
	}

	if err := kdController.Run(1, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
