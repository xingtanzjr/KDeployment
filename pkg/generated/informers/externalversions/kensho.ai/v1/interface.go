/*
This is the test project from Kensho
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "kensho.ai/kdeployment/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// KDeployments returns a KDeploymentInformer.
	KDeployments() KDeploymentInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// KDeployments returns a KDeploymentInformer.
func (v *version) KDeployments() KDeploymentInformer {
	return &kDeploymentInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
