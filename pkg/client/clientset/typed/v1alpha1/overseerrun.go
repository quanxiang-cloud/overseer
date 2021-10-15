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

type OverseerRunGetter interface {
	OverseerRuns(namespace string) OverseerRunsInterface
}

type OverseerRunsInterface interface {
	Create(ctx context.Context, overseer *v1alpha1.OverseerRun, opts v1.CreateOptions) (*v1alpha1.OverseerRun, error)
	Update(ctx context.Context, overseer *v1alpha1.OverseerRun, opts v1.UpdateOptions) (*v1alpha1.OverseerRun, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.OverseerRun, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.OverseerRunList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.OverseerRun, err error)
}

type overseerRun struct {
	client rest.Interface
	ns     string
}

func newOverseerRuns(c *OverseerV1alpha1Client, namespace string) *overseerRun {
	return &overseerRun{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

func (o *overseerRun) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.OverseerRun, err error) {
	result = &v1alpha1.OverseerRun{}
	err = o.client.Get().
		Namespace(o.ns).
		Resource("overseerruns").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

func (o *overseerRun) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.OverseerRunList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.OverseerRunList{}
	err = o.client.Get().
		Namespace(o.ns).
		Resource("overseerruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

func (o *overseerRun) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return o.client.Get().
		Namespace(o.ns).
		Resource("overseerruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (o *overseerRun) Create(ctx context.Context, overseerrun *v1alpha1.OverseerRun, opts v1.CreateOptions) (result *v1alpha1.OverseerRun, err error) {
	result = &v1alpha1.OverseerRun{}
	err = o.client.Post().
		Namespace(o.ns).
		Resource("overseerruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(overseerrun).
		Do(ctx).
		Into(result)
	return
}

func (o *overseerRun) Update(ctx context.Context, overseerrun *v1alpha1.OverseerRun, opts v1.UpdateOptions) (result *v1alpha1.OverseerRun, err error) {
	result = &v1alpha1.OverseerRun{}
	err = o.client.Put().
		Namespace(o.ns).
		Resource("overseerruns").
		Name(overseerrun.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(overseerrun).
		Do(ctx).
		Into(result)
	return
}

func (o *overseerRun) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return o.client.Delete().
		Namespace(o.ns).
		Resource("overseerruns").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (o *overseerRun) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return o.client.Delete().
		Namespace(o.ns).
		Resource("overseerruns").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

func (o *overseerRun) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.OverseerRun, err error) {
	result = &v1alpha1.OverseerRun{}
	err = o.client.Patch(pt).
		Namespace(o.ns).
		Resource("overseerruns").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
