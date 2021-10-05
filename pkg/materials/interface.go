package materials

import "github.com/quanxiang-cloud/overseer/pkg/materials/v1alpha1"

type Interface interface {
	V1alpha1() v1alpha1.Interface
}

type impl struct {
}

func New() Interface {
	return &impl{}
}

func (*impl) V1alpha1() v1alpha1.Interface {
	return v1alpha1.New()
}
