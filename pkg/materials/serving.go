package materials

import (
	"context"
	"time"

	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Serving struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *Serving) Reconcile(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, servingSpec *v1alpha1.ServingSpec) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcileServing")

	status := v1alpha1.VersatileStatus{}

	serving := &v1alpha1.Serving{}
	name := genName(osr.Name, vst.Name)
	err := r.Get(ctx, types.NamespacedName{Namespace: osr.Namespace, Name: name}, serving)
	if errors.IsNotFound(err) {
		serving, err = r.createServing(ctx, osr, vst, servingSpec)
		if err != nil {
			return err
		}

		status.Ref = serving.Name
		status.StartTime = metav1.NewTime(time.Now())
		status.Status = corev1.ConditionUnknown
	}

	if err != nil {
		logger.Error(err, "get serving", "name", name)
		return err
	}

	cond := serving.GetCondition(v1alpha1.Succeeded)
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

	status.Status = serving.Status.Status
	if status.Status == corev1.ConditionFalse {
		osr.Status.SetFalse()
	}

	osr.Status.SetVersatileStatus(status)
	return nil
}
func (r *Serving) createServing(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, servingSpec *v1alpha1.ServingSpec) (*v1alpha1.Serving, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createServing")

	serving := &v1alpha1.Serving{
		ObjectMeta: metav1.ObjectMeta{
			Name:      genName(osr.Name, vst.Name),
			Namespace: osr.Namespace,
		},
		Spec: *servingSpec.DeepCopy(),
	}

	if serving.Spec.ServiceAccountName == nil {
		serving.Spec.ServiceAccountName = &osr.Spec.ServiceAccountName
	}

	serving.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(osr, serving, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", serving.Name)
		return nil, err
	}

	if err := r.Create(ctx, serving); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}
	return serving, nil
}
