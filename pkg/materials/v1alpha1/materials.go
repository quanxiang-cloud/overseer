package v1alpha1

import (
	"fmt"

	"github.com/ghodss/yaml"
	apiv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	artifactsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/artifacts/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type materialsv1alpha1 struct {
	body      []byte
	namespace string
}

func New() *materialsv1alpha1 {
	return &materialsv1alpha1{}
}

func (m *materialsv1alpha1) Namespace(ns string) Interface {
	m.namespace = ns
	return m
}

func (m *materialsv1alpha1) Body(body []byte) Interface {
	m.body = body
	return m
}

func (m *materialsv1alpha1) Do() (client.Object, error) {
	if m.body == nil {
		return nil, nil
	}

	typeMeta := &struct {
		APIVersion string `yaml:"apiVersion,omitempty"`
		Kind       string `yaml:"kind,omitempty"`
	}{}
	err := yaml.Unmarshal(m.body, typeMeta)
	if err != nil {
		return nil, err
	}

	gvk := schema.FromAPIVersionAndKind(typeMeta.APIVersion, typeMeta.Kind)

	obj, ok := artifactsv1alpha1.GetObj(gvk)
	if !ok {
		return nil, fmt.Errorf("unrecognizable gkv [%v]", obj)
	}

	err = m.unmarshal(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *materialsv1alpha1) Param(params []apiv1alpha1.ParamSpec) Interface {
	// e.body = execute(paramsSliceToMap(params), string(e.body))
	return m
}

func (m *materialsv1alpha1) unmarshal(ptr interface{}) error {
	return yaml.Unmarshal(m.body, ptr)
}
