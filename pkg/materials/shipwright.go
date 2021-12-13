package materials

import (
	"context"
	"time"

	"github.com/quanxiang-cloud/overseer/pkg/apis/overseer/v1alpha1"
	shipwrightv1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Shipwright struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=shipwright.io,resources=builders;buildruns;buildstrategies;clusterbuildstrategies,verbs=get;list;watch;create;update;patch;delete

func (r *Shipwright) Reconcile(ctx context.Context, builder *v1alpha1.Builder) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcileShipwright")

	buildrun := &shipwrightv1alpha1.BuildRun{}

	if builder.Status.Ref == "" {
		// need to create builder
		var err error
		buildrun, err := r.create(ctx, builder)
		if err != nil {
			logger.Error(err, "name", builder.Name)
			return err
		}

		builder.Status.Ref = buildrun.Name
		builder.Status.StartTime = metav1.Time{Time: time.Now()}
		builder.Status.Status = corev1.ConditionUnknown
	} else {
		err := r.Get(ctx, types.NamespacedName{Namespace: builder.Namespace, Name: builder.Status.Ref}, buildrun)
		if err != nil {
			if errors.IsNotFound(err) {
				builder.Status.CompletionTime = metav1.Time{Time: time.Now()}
				builder.Status.Status = corev1.ConditionFalse
				builder.Status.Conditions = v1alpha1.Conditions{{
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

	cond := buildrun.Status.GetCondition(shipwrightv1alpha1.Succeeded)
	if cond != nil {
		builder.Status.CompletionTime = metav1.Time{Time: time.Now()}
		builder.Status.Status = cond.Status
		builder.Status.Conditions = v1alpha1.Conditions{{
			Type:               v1alpha1.Type(cond.Type),
			Status:             cond.Status,
			LastTransitionTime: cond.LastTransitionTime,
			Reason:             cond.Reason,
			Message:            cond.Message,
		}}
	}

	return nil
}

func (r *Shipwright) create(ctx context.Context, builder *v1alpha1.Builder) (*shipwrightv1alpha1.BuildRun, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createShipwright")

	build := &shipwrightv1alpha1.Build{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: builder.Name + "-",
			Namespace:    builder.Namespace,
		},
		Spec: shipwrightv1alpha1.BuildSpec{
			Source: shipwrightv1alpha1.Source{
				URL:        builder.Spec.Git.Url,
				Revision:   builder.Spec.Git.Revision,
				ContextDir: builder.Spec.Git.Subpath,
			},
			Strategy: &shipwrightv1alpha1.Strategy{
				Name: builder.Spec.Shipwright.Strategy.Name,
				Kind: (*shipwrightv1alpha1.BuildStrategyKind)(&builder.Spec.Shipwright.Strategy.Kind),
			},
			Output: shipwrightv1alpha1.Image{
				Image:       builder.Spec.Image.Image,
				Credentials: builder.Spec.Image.Credentials,
			},
		},
	}

	if length := len(builder.Spec.Params); length != 0 {
		build.Spec.ParamValues = make([]shipwrightv1alpha1.ParamValue, 0, length)
		for _, param := range builder.Spec.Params {
			build.Spec.ParamValues = append(build.Spec.ParamValues, shipwrightv1alpha1.ParamValue{
				Name:  param.Name,
				Value: param.Value,
			})
		}
	}

	build.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(builder, build, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", build.Name)
		return nil, err
	}
	if err := r.Create(ctx, build); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	buildrun := &shipwrightv1alpha1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      build.Name,
			Namespace: builder.Namespace,
		},
		Spec: shipwrightv1alpha1.BuildRunSpec{
			BuildRef: &shipwrightv1alpha1.BuildRef{
				Name: build.Name,
			},
			ServiceAccount: &shipwrightv1alpha1.ServiceAccount{
				Name: builder.Spec.ServiceAccountName,
			},
		},
	}

	buildrun.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(builder, buildrun, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", build.Name)
		return nil, err
	}
	if err := r.Create(ctx, buildrun); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	return buildrun, nil
}
