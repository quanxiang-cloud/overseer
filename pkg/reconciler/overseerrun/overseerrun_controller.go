/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package overseerrun

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	artifactsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/artifacts/v1alpha1"
	overseerRunV1alpha1 "github.com/quanxiang-cloud/overseer/pkg/listers/v1alpha1"
	"github.com/quanxiang-cloud/overseer/pkg/materials"
	materialsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/materials/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OverseerRunReconciler reconciles a OverseerRun object
type OverseerRunReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Log       logr.Logger
	lister    overseerRunV1alpha1.OverseerLister
	materials materials.Interface
}

func NewOverseerRunReconciler(client client.Client, scheme *runtime.Scheme, logger logr.Logger, lister overseerRunV1alpha1.OverseerLister) *OverseerRunReconciler {
	return &OverseerRunReconciler{
		Client:    client,
		Scheme:    scheme,
		Log:       logger.WithName("OverseerRun"),
		lister:    lister,
		materials: materials.New(),
	}
}

//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseerruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseerruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseerruns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OverseerRun object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *OverseerRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("Reconcile")

	var osr v1alpha1.OverseerRun
	err := r.Get(ctx, req.NamespacedName, &osr)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if osr.Status.Condition.IsTrue() || osr.Status.Condition.IsFalse() {
		return ctrl.Result{}, err
	}

	// phase is not done and the ref is finish,do next step
	if !osr.Status.Phase.IsDone() &&
		osr.Status.Condition.IsFinish(osr.Status.Phase.Sting()) {
		err = r.reconcile(ctx, &osr)
		if err != nil {
			failedWithError(&osr.Status.Status, err)
		}
	} else {
		err = r.updateStatus(ctx, &osr)
		if err != nil {
			failedWithError(&osr.Status.Status, err)
		}

	}

	if err = r.Status().Update(ctx, &osr); err != nil {
		logger.V(1).Error(err, "update status")
		return ctrl.Result{}, err
	}

	logger.Info("Success")
	return ctrl.Result{
		RequeueAfter: time.Second * 10,
	}, err
}

func (r *OverseerRunReconciler) updateStatus(ctx context.Context, osr *v1alpha1.OverseerRun) error {
	log := r.Log.WithName("updateStatus")

	if osr.Status.Phase.IsDone() {
		osr.Status.Condition.Status = corev1.ConditionTrue
		return nil
	}

	name := osr.Status.Phase.Sting()
	ref := osr.Status.Condition.ResourceRef[name]

	getter, ok := artifactsv1alpha1.GetGetter(ref.GroupVersionKind)
	if !ok {
		err := fmt.Errorf("unknown gvk [%s]", ref.GroupVersionKind)
		log.Error(err, "get object by gvk", "refName", name)
		ref.State = v1alpha1.StepConditionFail
		return err
	}

	obj := getter.New()

	if err := r.Get(ctx, client.ObjectKey{Namespace: osr.Namespace, Name: ref.RefName}, obj); err != nil {
		log.Error(err, "failed to get object", "refName", ref.RefName)
		ref.State = v1alpha1.StepConditionFail
		return err
	}

	sc := getter.GetState(obj)
	sc.RefName = ref.RefName

	if sc.State == v1alpha1.StepConditionFail {
		osr.Status.Condition.Status = corev1.ConditionFalse
	}
	osr.Status.Condition.ResourceRef[name] = sc
	osr.Status.Condition.LastTransitionTime = metav1.NewTime(time.Now())
	return nil
}

func (r *OverseerRunReconciler) reconcile(ctx context.Context, osr *v1alpha1.OverseerRun) error {
	overseer, err := r.lister.Overseers(osr.Namespace).Get(osr.Spec.OverseerRef.Name)
	if err != nil {
		return err
	}

	err = v1alpha1.ParamsValidate(osr.Spec.Params, overseer.Spec.Params)
	if err != nil {
		return err
	}

	if osr.Status.Condition.IsNil() {
		osr.Status.Condition.Init()
	}

	var getNext = func(overseer *v1alpha1.Overseer, osr *v1alpha1.OverseerRun) *v1alpha1.StepSpec {
		if len(overseer.Spec.Steps) == 0 {
			return nil
		}
		if osr.Status.Phase.IsNil() {
			return overseer.Spec.Steps[0].DeepCopy()
		}

		var index int
		for i, step := range overseer.Spec.Steps {
			if osr.Status.Phase.Equal(step.Name) {
				index = i
				break
			}
		}

		if index+1 < len(overseer.Spec.Steps) {
			return overseer.Spec.Steps[index+1].DeepCopy()
		}

		return nil
	}

	next := getNext(overseer, osr)
	if next == nil {
		osr.Status.Phase = v1alpha1.PhaseDone
		return nil
	}

	err = r.reconcileStep(ctx, next, osr)
	if err != nil {
		return err
	}

	return nil
}

func (r *OverseerRunReconciler) reconcileStep(ctx context.Context, step *v1alpha1.StepSpec, osr *v1alpha1.OverseerRun) error {
	log := r.Log.WithName("reconcileStep")

	m := r.materials.V1alpha1()
	obj, err := m.
		Body([]byte(step.Template)).
		Param(osr.Spec.Params).
		Do(materialsv1alpha1.WithNamespace(osr.Namespace),
			materialsv1alpha1.WithAttachedGenerateName(osr.Name),
		)

	if err != nil {
		return err
	}

	obj.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(osr, obj, r.Scheme); err != nil {
		log.Error(err, "Failed to SetControllerReference", "step", step.Name)
		return err
	}

	if err = r.Create(ctx, obj); err != nil {
		log.Error(err, "Failed to create obj")
		return err
	}

	osr.Status.Phase = v1alpha1.Phase(step.Name)
	osr.Status.Condition.ResourceRef[step.Name] = v1alpha1.StepCondition{
		GroupVersionKind: m.GetGroupVersionKind().String(),
		RefName:          obj.GetName(),
		State:            v1alpha1.StepConditionUnknown,
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OverseerRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OverseerRun{}).
		Complete(r)
}
