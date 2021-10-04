package v1alpha1

import (
	"context"
	"time"

	"github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	"github.com/quanxiang-cloud/overseer/pkg/client/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type OverseerGetter interface {
	Overseers(namespace string) OverseerInterface
}

type OverseerInterface interface {
	Create(ctx context.Context, overseer *v1alpha1.Overseer, opts v1.CreateOptions) (*v1alpha1.Overseer, error)
	Update(ctx context.Context, overseer *v1alpha1.Overseer, opts v1.UpdateOptions) (*v1alpha1.Overseer, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Overseer, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.OverseerList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Overseer, err error)
}

type overseer struct {
	client rest.Interface
	ns     string
}

func newOverseers(c *OverseerV1alpha1Client, namespace string) *overseer {
	return &overseer{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

func (o *overseer) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Overseer, err error) {
	result = &v1alpha1.Overseer{}
	err = o.client.Get().
		Namespace(o.ns).
		Resource("overseers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (o *overseer) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.OverseerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.OverseerList{}
	err = o.client.Get().
		Namespace(o.ns).
		Resource("overseers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

func (o *overseer) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return o.client.Get().
		Namespace(o.ns).
		Resource("overseers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (o *overseer) Create(ctx context.Context, pipeline *v1alpha1.Overseer, opts v1.CreateOptions) (result *v1alpha1.Overseer, err error) {
	result = &v1alpha1.Overseer{}
	err = o.client.Post().
		Namespace(o.ns).
		Resource("overseers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(pipeline).
		Do(ctx).
		Into(result)
	return
}

func (o *overseer) Update(ctx context.Context, pipeline *v1alpha1.Overseer, opts v1.UpdateOptions) (result *v1alpha1.Overseer, err error) {
	result = &v1alpha1.Overseer{}
	err = o.client.Put().
		Namespace(o.ns).
		Resource("overseers").
		Name(pipeline.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(pipeline).
		Do(ctx).
		Into(result)
	return
}

func (o *overseer) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return o.client.Delete().
		Namespace(o.ns).
		Resource("overseers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (o *overseer) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return o.client.Delete().
		Namespace(o.ns).
		Resource("overseers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (o *overseer) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Overseer, err error) {
	result = &v1alpha1.Overseer{}
	err = o.client.Patch(pt).
		Namespace(o.ns).
		Resource("overseers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
