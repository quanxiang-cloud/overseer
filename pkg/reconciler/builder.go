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

package reconciler

import (
	"context"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/apis/overseer/v1alpha1"
	"github.com/quanxiang-cloud/overseer/pkg/materials"
)

// BuilderReconciler reconciles a Builder object
type BuilderReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	*materials.Shipwright
}

//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=builders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=builders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=builders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Builder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *BuilderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("Reconcile")

	var builder v1alpha1.Builder
	err := r.Get(ctx, req.NamespacedName, &builder)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if builder.IsDone() {
		return ctrl.Result{}, nil
	}

	status := builder.Status.DeepCopy()

	if builder.Spec.Shipwright != nil {
		err := r.Shipwright.Reconcile(ctx, &builder)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if reflect.DeepEqual(builder.Status, status) {
		return ctrl.Result{
			RequeueAfter: time.Second * 10,
		}, nil
	}

	if err = r.Status().Update(ctx, &builder); err != nil {
		logger.Error(err, "update status")
		return ctrl.Result{}, err
	}

	logger.Info("Success")
	return ctrl.Result{
		RequeueAfter: time.Second * 10,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuilderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Builder{}).
		Complete(r)
}
