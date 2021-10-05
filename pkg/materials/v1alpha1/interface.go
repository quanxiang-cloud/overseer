package v1alpha1

import (
	apiv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	artifactsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/artifacts/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface interface {
	Body(body []byte) Interface
	Param(params []apiv1alpha1.ParamSpec) Interface
	Do(...artifactsv1alpha1.Options) (client.Object, error)
}
