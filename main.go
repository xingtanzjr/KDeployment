package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	clientset "kensho.ai/kdeployment/pkg/generated/clientset/versioned"
)

func main() {
	klog.InitFlags(nil)

	kubeconfig := "/Users/zhangjinrui/.kube/config"

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Error("Cannot parse kubeconfig")
	}

	// defaultClient, err := kubernetes.NewForConfig(cfg)
	kdeploymentClient, err := clientset.NewForConfig(cfg)

	kd, err := kdeploymentClient.DistributionV1().KDeployments("default").Get(context.TODO(), "kdeployment-sample", metav1.GetOptions{})
	if err != nil {
		klog.Error(err, "Error when read")
	}
	klog.Info(kd.Name)
}
