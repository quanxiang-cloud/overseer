package clientset

import (
	"fmt"

	"github.com/quanxiang-cloud/overseer/pkg/client/clientset/typed/v1alpha1"
	"k8s.io/client-go/discovery"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	V1alpha1() v1alpha1.Interface
}

type Clientset struct {
	*kubernetes.Clientset
	v1alpha1 *v1alpha1.OverseerV1alpha1Client
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}
func (c *Clientset) V1alpha1() v1alpha1.Interface {
	return c.v1alpha1
}

func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.v1alpha1, err = v1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.Clientset, err = kubernetes.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.v1alpha1 = v1alpha1.NewForConfigOrDie(c)

	cs.Clientset = kubernetes.NewForConfigOrDie(c)
	return &cs
}

func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.v1alpha1 = v1alpha1.New(c)

	cs.Clientset = kubernetes.New(c)
	return &cs
}
