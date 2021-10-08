package v1alpha1

import (
	"github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	"github.com/quanxiang-cloud/overseer/pkg/client/clientset/scheme"
	"k8s.io/client-go/rest"
)

type Interface interface {
	RESTClient() rest.Interface
	OverseerGetter
	OverseerRunGetter
}

type OverseerV1alpha1Client struct {
	restClient rest.Interface
}

func (o *OverseerV1alpha1Client) Overseers(namespace string) OverseerInterface {
	return newOverseers(o, namespace)
}

func (o *OverseerV1alpha1Client) OverseerRuns(namespace string) OverseerRunsInterface {
	return newOverseerRuns(o, namespace)
}

func (o *OverseerV1alpha1Client) RESTClient() rest.Interface {
	return o.restClient
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.GroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

func NewForConfig(c *rest.Config) (*OverseerV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &OverseerV1alpha1Client{client}, nil
}

func NewForConfigOrDie(c *rest.Config) *OverseerV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

func New(c rest.Interface) *OverseerV1alpha1Client {
	return &OverseerV1alpha1Client{c}
}
