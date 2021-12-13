package materials

import (
	"context"
	"time"

	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/apis/overseer/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Builder struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *Builder) Reconcile(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, builderSpec *v1alpha1.BuilderSpec) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcileBuilder")

	status := v1alpha1.VersatileStatus{}

	builder := &v1alpha1.Builder{}
	name := genName(osr.Name, vst.Name)
	err := r.Get(ctx, types.NamespacedName{Namespace: osr.Namespace, Name: name}, builder)
	if errors.IsNotFound(err) {
		builder, err = r.createBuilder(ctx, osr, vst, builderSpec)
		if err != nil {
			return err
		}
		status.Ref = builder.Name
		status.StartTime = metav1.NewTime(time.Now())
		status.Status = corev1.ConditionUnknown
	}

	if err != nil {
		logger.Error(err, "get builder", "name", name)
		return err
	}

	cond := builder.GetCondition(v1alpha1.Succeeded)
	if cond != nil {
		status.Conditions = v1alpha1.Conditions{
			{
				Type:               v1alpha1.Succeeded,
				Status:             cond.Status,
				LastTransitionTime: cond.LastTransitionTime,
				Reason:             cond.Reason,
				Message:            cond.Message,
			},
		}
	}

	status.Status = builder.Status.Status
	if status.Status == corev1.ConditionFalse {
		osr.Status.SetFalse()
	}

	osr.Status.SetVersatileStatus(status)
	return nil
}

func (r *Builder) createBuilder(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, builderSpec *v1alpha1.BuilderSpec) (*v1alpha1.Builder, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createBuilder")

	builder := &v1alpha1.Builder{
		ObjectMeta: metav1.ObjectMeta{
			Name:      genName(osr.Name, vst.Name),
			Namespace: osr.Namespace,
			Labels:    osr.Labels,
		},
		Spec: *builderSpec.DeepCopy(),
	}

	if builder.Spec.ServiceAccountName == nil {
		builder.Spec.ServiceAccountName = &osr.Spec.ServiceAccountName
	}

	builder.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(osr, builder, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", builder.Name)
		return nil, err
	}

	if err := r.Create(ctx, builder); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	return builder, nil
}
