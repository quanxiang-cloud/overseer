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

// ServingReconciler reconciles a Serving object
type ServingReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	*materials.Knative
	*materials.Service
}

//+kubebuilder:rbac:groups=overseer.quanxiang.cloud.io,resources=servings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=overseer.quanxiang.cloud.io,resources=servings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=overseer.quanxiang.cloud.io,resources=servings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Serving object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ServingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("Reconcile")

	var serving v1alpha1.Serving
	err := r.Get(ctx, req.NamespacedName, &serving)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if serving.IsDone() {
		return ctrl.Result{}, nil
	}

	status := serving.Status.DeepCopy()

	if serving.Spec.Knative != nil {
		err = r.Knative.Reconcile(ctx, &serving)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if serving.Spec.Service != nil {
		err = r.Service.Reconcile(ctx, &serving)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if reflect.DeepEqual(serving.Status, status) {
		return ctrl.Result{
			RequeueAfter: time.Second * 10,
		}, nil
	}

	if err = r.Status().Update(ctx, &serving); err != nil {
		logger.Error(err, "update status")
		return ctrl.Result{}, err
	}

	logger.Info("Success")
	return ctrl.Result{
		RequeueAfter: time.Second * 10,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Serving{}).
		Complete(r)
}
