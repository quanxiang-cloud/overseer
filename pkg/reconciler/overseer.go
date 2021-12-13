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
	"fmt"
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

// OverseerReconciler reconciles a Overseer object
type OverseerReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	*materials.PipelineRun
	*materials.Builder
	*materials.Serving
}

//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=quanxiang.cloud.io,resources=overseers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Overseer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *OverseerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("Reconcile")

	var osr v1alpha1.Overseer
	err := r.Get(ctx, req.NamespacedName, &osr)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if osr.IsDone() {
		return ctrl.Result{}, nil
	}

	status := osr.Status.DeepCopy()

	err = r.reconcile(ctx, &osr)
	if err != nil {
		return ctrl.Result{}, err
	}

	if reflect.DeepEqual(osr.Status, status) {
		return ctrl.Result{
			RequeueAfter: time.Second * 10,
		}, nil
	}

	if err = r.Status().Update(ctx, &osr); err != nil {
		logger.Error(err, "update status")
		return ctrl.Result{}, err
	}

	logger.Info("Success")
	return ctrl.Result{
		RequeueAfter: time.Second * 10,
	}, nil
}

func (r *OverseerReconciler) reconcile(ctx context.Context, osr *v1alpha1.Overseer) error {
	// get the next task when the current step is completed.
	// end the task without the next step and mark success.
	if !osr.ShoudContinue() {
		osr.Status.SetSuccess()
		return nil
	}

	var err error

	vst := osr.GetVersatile()
	switch osr.Status.Phase.Stage {
	case v1alpha1.PipelineRunStage:
		err = r.PipelineRun.Reconcile(ctx, osr, vst, vst.PipelineRun)
	case v1alpha1.BuilderStage:
		err = r.Builder.Reconcile(ctx, osr, vst, vst.Builder)
	case v1alpha1.ServingStage:
		err = r.Serving.Reconcile(ctx, osr, vst, vst.Serving)
	default:
		return fmt.Errorf("unexpected execution phase")
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *OverseerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Overseer{}).
		Complete(r)
}
