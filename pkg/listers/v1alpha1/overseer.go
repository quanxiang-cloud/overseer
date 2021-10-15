package v1alpha1

import (
	"github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type OverseerLister interface {
	List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error)
	Overseers(namespace string) OverseerNamespaceLister
}

type overseeLister struct {
	indexer cache.Indexer
}

func NewOverseerLister(indexer cache.Indexer) OverseerLister {
	return &overseeLister{indexer: indexer}
}

func (s *overseeLister) List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Overseer))
	})
	return ret, err
}

func (s *overseeLister) Overseers(namespace string) OverseerNamespaceLister {
	return overseerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// OverseerNamespaceLister helps list Overseers.
type OverseerNamespaceLister interface {
	List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error)
	Get(name string) (*v1alpha1.Overseer, error)
}

// overseerLister implements the OverseerLister interface.
type overseerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

func (o overseerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error) {
	err = cache.ListAllByNamespace(o.indexer, o.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Overseer))
	})
	return ret, err
}

func (o overseerNamespaceLister) Get(name string) (*v1alpha1.Overseer, error) {
	obj, exists, err := o.indexer.GetByKey(o.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("overseer"), name)
	}
	return obj.(*v1alpha1.Overseer), nil
}
