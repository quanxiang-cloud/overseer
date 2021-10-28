package v1alpha1

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ghodss/yaml"
	apiv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	artifactsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/artifacts/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type materialsv1alpha1 struct {
	body []byte

	params []apiv1alpha1.Param
	gvk    schema.GroupVersionKind
}

func New() *materialsv1alpha1 {
	return &materialsv1alpha1{}
}

func (m *materialsv1alpha1) Body(body []byte) Interface {
	m.body = body
	return m
}

func (m *materialsv1alpha1) Do(opts ...Options) (client.Object, error) {
	if m.body == nil {
		return nil, nil
	}

	err := m.param()
	if err != nil {
		return nil, err
	}

	typeMeta := &struct {
		APIVersion string `yaml:"apiVersion,omitempty"`
		Kind       string `yaml:"kind,omitempty"`
	}{}

	err = yaml.Unmarshal(m.body, typeMeta)
	if err != nil {
		return nil, err
	}

	m.gvk = schema.FromAPIVersionAndKind(typeMeta.APIVersion, typeMeta.Kind)

	obj, ok := artifactsv1alpha1.GetObj(m.gvk)
	if !ok {
		return nil, fmt.Errorf("unrecognizable gkv [%s/%s]", typeMeta.APIVersion, typeMeta.Kind)
	}

	err = m.unmarshal(obj)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(obj)
	}
	return obj, nil
}

func (m *materialsv1alpha1) Param(params []apiv1alpha1.Param) Interface {
	m.params = params
	return m
}

func (m *materialsv1alpha1) param() error {
	kv := make(map[string]string, len(m.params))
	for _, param := range m.params {
		kv[param.Name] = param.Value
	}

	body, err := templateExecute(kv, string(m.body))
	if err != nil {
		return err
	}

	m.body = []byte(body)
	return nil
}

func templateExecute(values interface{}, tmpl string) (string, error) {
	var err error
	buf := &bytes.Buffer{}

	tf := template.New("")
	tf.Delims("$(params", ")")

	tf, err = tf.Parse(string(tmpl))
	if err != nil {
		return tmpl, err
	}

	err = tf.Execute(buf, values)
	if err != nil {
		return tmpl, err
	}
	fmt.Println(buf.String())
	return buf.String(), nil
}

func (m *materialsv1alpha1) unmarshal(ptr interface{}) error {
	return yaml.Unmarshal(m.body, ptr)
}

func (m *materialsv1alpha1) GetGroupVersionKind() schema.GroupVersionKind {
	return m.gvk
}
