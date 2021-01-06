/*
This is the test project from Kensho
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1 "kensho.ai/kdeployment/pkg/apis/distribution.kensho.ai/v1"
)

// KDeploymentLister helps list KDeployments.
// All objects returned here must be treated as read-only.
type KDeploymentLister interface {
	// List lists all KDeployments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.KDeployment, err error)
	// KDeployments returns an object that can list and get KDeployments.
	KDeployments(namespace string) KDeploymentNamespaceLister
	KDeploymentListerExpansion
}

// kDeploymentLister implements the KDeploymentLister interface.
type kDeploymentLister struct {
	indexer cache.Indexer
}

// NewKDeploymentLister returns a new KDeploymentLister.
func NewKDeploymentLister(indexer cache.Indexer) KDeploymentLister {
	return &kDeploymentLister{indexer: indexer}
}

// List lists all KDeployments in the indexer.
func (s *kDeploymentLister) List(selector labels.Selector) (ret []*v1.KDeployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.KDeployment))
	})
	return ret, err
}

// KDeployments returns an object that can list and get KDeployments.
func (s *kDeploymentLister) KDeployments(namespace string) KDeploymentNamespaceLister {
	return kDeploymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KDeploymentNamespaceLister helps list and get KDeployments.
// All objects returned here must be treated as read-only.
type KDeploymentNamespaceLister interface {
	// List lists all KDeployments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.KDeployment, err error)
	// Get retrieves the KDeployment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.KDeployment, error)
	KDeploymentNamespaceListerExpansion
}

// kDeploymentNamespaceLister implements the KDeploymentNamespaceLister
// interface.
type kDeploymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all KDeployments in the indexer for a given namespace.
func (s kDeploymentNamespaceLister) List(selector labels.Selector) (ret []*v1.KDeployment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.KDeployment))
	})
	return ret, err
}

// Get retrieves the KDeployment from the indexer for a given namespace and name.
func (s kDeploymentNamespaceLister) Get(name string) (*v1.KDeployment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("kdeployment"), name)
	}
	return obj.(*v1.KDeployment), nil
}