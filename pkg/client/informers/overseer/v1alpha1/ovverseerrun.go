package v1alpha1

import (
	"context"
	"time"

	overseerv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	clientv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/client/clientset"
	"github.com/quanxiang-cloud/overseer/pkg/client/informers"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/tools/cache"
)

// OverseerRunInformer provides access to a shared informer and lister
type OverseerRunInformer interface {
	Informer() cache.SharedIndexInformer
}

func NewOverseerRunInformer(client clientv1alpha1.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.V1alpha1().OverseerRuns(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.V1alpha1().OverseerRuns(namespace).Watch(context.TODO(), options)
			},
		},
		&overseerv1alpha1.OverseerRun{},
		resyncPeriod,
		indexers,
	)
}

type overseerRunInformer struct {
	factory          informers.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

func (o *overseerRunInformer) defaultInformer(client clientv1alpha1.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewOverseerRunInformer(client, o.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, o.tweakListOptions)
}

func (o *overseerRunInformer) Informer() cache.SharedIndexInformer {
	return o.factory.InformerFor(&overseerv1alpha1.OverseerRun{}, o.defaultInformer)
}
