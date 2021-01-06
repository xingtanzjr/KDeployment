package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	v1 "kensho.ai/kdeployment/pkg/apis/distribution.kensho.ai/v1"
	clientset "kensho.ai/kdeployment/pkg/generated/clientset/versioned"
	kdInformers "kensho.ai/kdeployment/pkg/generated/informers/externalversions/distribution.kensho.ai/v1"
	kdlisters "kensho.ai/kdeployment/pkg/generated/listers/distribution.kensho.ai/v1"
)

type KDEventItem struct {
	clusterId string
	key       string
}

type KDClusterTool struct {
	clusterId string
	lister    kdlisters.KDeploymentLister
	client    clientset.Interface
	synced    cache.InformerSynced
}

type KDController struct {
	clusterToolMap map[string]KDClusterTool
	workqueue      workqueue.RateLimitingInterface
	//TODO to learn the recorder
}

func NewKDController(informerMap map[string]kdInformers.KDeploymentInformer, clientMap map[string]clientset.Interface) *KDController {
	clusterToolMap := make(map[string]KDClusterTool)

	for k, informer := range informerMap {
		tool := KDClusterTool{
			clusterId: k,
			lister:    informer.Lister(),
			synced:    informer.Informer().HasSynced,
			client:    clientMap[k],
		}
		clusterToolMap[k] = tool
	}

	kdController := &KDController{
		clusterToolMap: clusterToolMap,
		workqueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "KDeployment"),
	}

	// set up event handler for each Informer
	for clusterId, informer := range informerMap {
		informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				kdController.enqueueKDeployment(clusterId, obj)
			},
			UpdateFunc: func(old, new interface{}) {
				kdController.enqueueKDeployment(clusterId, new)
			},
		})
	}

	return kdController
}

func (c *KDController) enqueueKDeployment(clusterId string, obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(KDEventItem{
		clusterId: clusterId,
		key:       key,
	})
}

func (c *KDController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting KDeployment controller")
	klog.Info("Waiting for informer caches to sync")
	for clusterId, clusterTool := range c.clusterToolMap {
		klog.Info("Waiting for " + clusterId)
		if ok := cache.WaitForCacheSync(stopCh, clusterTool.synced); !ok {
			return fmt.Errorf("failed to wait for caches to sycn for cluster %s", clusterId)
		}
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *KDController) runWorker() {

}

func (c *KDController) processNextKDEventItem() bool {
	kdEventItem, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var item KDEventItem
		var ok bool
		if item, ok = obj.(KDEventItem); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected KDEventItem in workqueue but got %#v", obj))
			return nil
		}
		err := c.processOneKDEventItem(item)
		if err != nil {
			c.workqueue.AddRateLimited(item)
			return fmt.Errorf("error when processing KDEventItem %s, %s", item.clusterId, item.key)
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully process KDEventItem %s, ^%s", item.clusterId, item.key)
		return nil
	}(kdEventItem)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *KDController) processOneKDEventItem(item KDEventItem) error {
	// get corresponding KDeployment
	namespace, name, err := cache.SplitMetaNamespaceKey(item.key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", item.key))
	}
	kdeployment, err := c.getKDeployment(item.clusterId, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("KDeployment '%s' in work queue no longer exists", item.key))
			return nil
		}
		return err
	}

	// Replicate KDeployment in each cluster
	for _, clustertool := range c.clusterToolMap {
		c.replicaKDeployment(clustertool, kdeployment)
	}

	return nil
}

// Ensure that the KDeployment will be replicated to each cluster
func (c *KDController) replicaKDeployment(clustertool KDClusterTool, kdeployment *v1.KDeployment) error {
	lister := clustertool.lister
	_, err := lister.KDeployments(kdeployment.Namespace).Get(kdeployment.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err := clustertool.client.DistributionV1().KDeployments(kdeployment.Namespace).Create(context.TODO(), kdeployment, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			klog.Info("create KDeployment[%s] on cluster [%s]", kdeployment.Name, clustertool.clusterId)
		} else {
			return err
		}
	}
	return nil
}

func (c *KDController) getKDeployment(clusterId string, namespace string, name string) (*v1.KDeployment, error) {
	kdeployment, err := c.clusterToolMap[clusterId].client.DistributionV1().KDeployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return kdeployment, err
}
