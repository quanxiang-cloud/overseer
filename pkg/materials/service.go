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

type Service struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *Service) Reconcile(ctx context.Context, serving *v1alpha1.Serving) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcileService")

	svc := &corev1.Service{}

	if serving.Status.Ref == "" {
		svc, err := r.create(ctx, serving)
		if err != nil {
			return err
		}

		serving.Status.Ref = svc.Name
		serving.Status.StartTime = metav1.Time{Time: time.Now()}
		serving.Status.Status = corev1.ConditionUnknown
	} else {
		err := r.Get(ctx, types.NamespacedName{Namespace: serving.Namespace, Name: serving.Status.Ref}, svc)
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

	serving.Status.CompletionTime = metav1.Time{Time: time.Now()}
	serving.Status.Status = corev1.ConditionTrue

	return nil
}

func (r *Service) create(ctx context.Context, serving *v1alpha1.Serving) (*corev1.Service, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createServing")

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: serving.Spec.Name + "-",
			Namespace:    serving.Namespace,
		},

		Spec: *serving.Spec.Template,
	}

	pod.Spec.ServiceAccountName = *serving.Spec.ServiceAccountName

	labels := pod.GetLabels()
	if len(labels) == 0 {
		labels = make(map[string]string)
	}
	labels["overseer.serving/name"] = serving.Name
	pod.SetLabels(labels)

	pod.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(serving, pod, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", pod.Name)
		return nil, err
	}
	if err := r.Create(ctx, pod); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serving.Spec.Name,
			Namespace: serving.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Ports:    serving.Spec.Service.Ports,
			Selector: labels,
			Type:     "ClusterIP",
		},
	}

	svc.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(serving, svc, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", svc.Name)
		return nil, err
	}
	if err := r.Create(ctx, svc); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	return svc, nil
}
