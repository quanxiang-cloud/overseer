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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	osv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	artifactsv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/artifacts/v1alpha1"
	overseerRunV1alpha1 "github.com/quanxiang-cloud/overseer/pkg/listers/v1alpha1"
	"github.com/quanxiang-cloud/overseer/pkg/materials"
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

	var osr osv1alpha1.OverseerRun
	err := r.Get(ctx, req.NamespacedName, &osr)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.reconcile(ctx, osr)
	if err != nil {
		return ctrl.Result{}, nil
	}

	// TODO osr	status
	return ctrl.Result{}, nil
}

func (r *OverseerRunReconciler) reconcile(ctx context.Context, osr osv1alpha1.OverseerRun) error {
	overseer, err := r.lister.Overseers(osr.Namespace).Get(osr.Spec.OverseerRef.Name)
	if err != nil {
		return err
	}

	err = osv1alpha1.ParamsValidate(osr.Spec.Params, overseer.Spec.Params)
	if err != nil {
		// TODO params error state
		return err
	}

	err = r.reconcileOverseer(ctx, osr, overseer)
	if err != nil {
		// TODO  should set status and reason
		return err
	}

	return nil
}

func (r *OverseerRunReconciler) reconcileOverseer(ctx context.Context, osr osv1alpha1.OverseerRun, overseer *osv1alpha1.Overseer) error {
	log := r.Log.WithName("reconcileOverseer")

	for _, step := range overseer.Spec.Steps {
		obj, err := r.materials.V1alpha1().
			Body([]byte(step.Template)).
			Param(osr.Spec.Params).
			Do(artifactsv1alpha1.WithNamespace(overseer.Namespace),
				artifactsv1alpha1.WithAttachedGenerateName(osr.Name),
			)

		if err != nil {
			return err
		}

		if err = r.Create(ctx, obj); err != nil {
			log.Error(err, "Failed to create obj")
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OverseerRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&osv1alpha1.OverseerRun{}).
		Complete(r)
}
