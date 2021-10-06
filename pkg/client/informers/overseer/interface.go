package overseer

import (
	"github.com/quanxiang-cloud/overseer/pkg/client/informers"
	"github.com/quanxiang-cloud/overseer/pkg/client/informers/overseer/v1alpha1"
	"k8s.io/client-go/informers/internalinterfaces"
)

type Interface interface {
	V1alpha1() v1alpha1.Interface
}

type impl struct {
	factory          informers.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func New(f informers.SharedInformerFactory, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &impl{factory: f, tweakListOptions: tweakListOptions}
}

func (i *impl) V1alpha1() v1alpha1.Interface {
	return v1alpha1.New(i.factory, i.tweakListOptions)
}
