package v1alpha1

import (
	"fmt"

	apiv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface interface {
	Body(body []byte) Interface
	Param(params []apiv1alpha1.Param) Interface
	Do(...Options) (client.Object, error)

	GetGroupVersionKind() schema.GroupVersionKind
}

type Options func(client.Object)

func WithNamespace(namespace string) Options {
	return func(obj client.Object) {
		obj.SetNamespace(namespace)
	}
}

func WithAttachedGenerateName(name string) Options {
	return func(obj client.Object) {
		obj.SetGenerateName(fmt.Sprintf("%s-%s-", name, obj.GetName()))
		obj.SetName("")
	}
}
