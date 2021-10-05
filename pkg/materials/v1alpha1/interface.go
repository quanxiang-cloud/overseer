package v1alpha1

import (
	apiv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface interface {
	Namespace(namespace string) Interface
	Body(body []byte) Interface
	Param(params []apiv1alpha1.ParamSpec) Interface
	Do() (client.Object, error)
}
