/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/apis/overseer/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// OverseerLister helps list Overseers.
type OverseerLister interface {
	// List lists all Overseers in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error)
	// Overseers returns an object that can list and get Overseers.
	Overseers(namespace string) OverseerNamespaceLister
	OverseerListerExpansion
}

// overseerLister implements the OverseerLister interface.
type overseerLister struct {
	indexer cache.Indexer
}

// NewOverseerLister returns a new OverseerLister.
func NewOverseerLister(indexer cache.Indexer) OverseerLister {
	return &overseerLister{indexer: indexer}
}

// List lists all Overseers in the indexer.
func (s *overseerLister) List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Overseer))
	})
	return ret, err
}

// Overseers returns an object that can list and get Overseers.
func (s *overseerLister) Overseers(namespace string) OverseerNamespaceLister {
	return overseerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// OverseerNamespaceLister helps list and get Overseers.
type OverseerNamespaceLister interface {
	// List lists all Overseers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error)
	// Get retrieves the Overseer from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Overseer, error)
	OverseerNamespaceListerExpansion
}

// overseerNamespaceLister implements the OverseerNamespaceLister
// interface.
type overseerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Overseers in the indexer for a given namespace.
func (s overseerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Overseer, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Overseer))
	})
	return ret, err
}

// Get retrieves the Overseer from the indexer for a given namespace and name.
func (s overseerNamespaceLister) Get(name string) (*v1alpha1.Overseer, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("overseer"), name)
	}
	return obj.(*v1alpha1.Overseer), nil
}
