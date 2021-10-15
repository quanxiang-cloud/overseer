package v1alpha1

import (
	"github.com/quanxiang-cloud/overseer/pkg/client/informers"
	"k8s.io/client-go/informers/internalinterfaces"
)

type Interface interface {
	Overseer() OverseerInformer
	OverseerRun() OverseerRunInformer
}

type informer struct {
	factory          informers.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func New(f informers.SharedInformerFactory, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &informer{factory: f, tweakListOptions: tweakListOptions}
}

func (i *informer) Overseer() OverseerInformer {
	return &overseerInformer{factory: i.factory, tweakListOptions: i.tweakListOptions}
}

func (i *informer) OverseerRun() OverseerRunInformer {
	return &overseerRunInformer{factory: i.factory, tweakListOptions: i.tweakListOptions}
}
