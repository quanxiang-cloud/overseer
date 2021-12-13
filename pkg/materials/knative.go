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
	knativeapi "knative.dev/pkg/apis"
	ksvcv1 "knative.dev/serving/pkg/apis/serving/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Knative struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *Knative) Reconcile(ctx context.Context, serving *v1alpha1.Serving) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcileKnative")

	ksvc := &ksvcv1.Service{}

	var err error
	if serving.Status.Ref == "" {
		ksvc, err = r.create(ctx, serving)
		if err != nil {
			return err
		}

		serving.Status.Ref = ksvc.Name
		serving.Status.StartTime = metav1.Time{Time: time.Now()}
		serving.Status.Status = corev1.ConditionUnknown
	} else {
		err := r.Get(ctx, types.NamespacedName{Namespace: serving.Namespace, Name: serving.Status.Ref}, ksvc)
		if err != nil {
			if errors.IsNotFound(err) {
				serving.Status.CompletionTime = metav1.Time{Time: time.Now()}
				serving.Status.Status = corev1.ConditionFalse
				serving.Status.Conditions = v1alpha1.Conditions{{
					Type:               v1alpha1.Succeeded,
					Status:             corev1.ConditionFalse,
					LastTransitionTime: metav1.Time{Time: time.Now()},
					Reason:             v1alpha1.Cancel,
				}}
				return nil
			}
			return err
		}
	}

	cond := ksvc.Status.GetCondition(knativeapi.ConditionSucceeded)
	if cond != nil {
		serving.Status.CompletionTime = metav1.Time{Time: time.Now()}
		serving.Status.Status = cond.Status
		serving.Status.Conditions = v1alpha1.Conditions{{
			Type:               v1alpha1.Type(cond.Type),
			Status:             cond.Status,
			LastTransitionTime: cond.LastTransitionTime.Inner,
			Reason:             cond.Reason,
			Message:            cond.Message,
		}}
	}

	return nil
}

func (r *Knative) create(ctx context.Context, serving *v1alpha1.Serving) (*ksvcv1.Service, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createKnative")

	ksvc := &ksvcv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serving.Spec.Name,
			Namespace: serving.Namespace,
		},
		Spec: ksvcv1.ServiceSpec{
			ConfigurationSpec: ksvcv1.ConfigurationSpec{
				Template: ksvcv1.RevisionTemplateSpec{
					Spec: ksvcv1.RevisionSpec{
						PodSpec: *serving.Spec.Template,
					},
				},
			},
		},
	}

	ksvc.Spec.Template.Spec.ServiceAccountName = *serving.Spec.ServiceAccountName

	ksvc.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(serving, ksvc, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", ksvc.Name)
		return nil, err
	}
	if err := r.Create(ctx, ksvc); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	return ksvc, nil
}
