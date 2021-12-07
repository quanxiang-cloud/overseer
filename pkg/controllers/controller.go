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

package controllers

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/quanxiang-cloud/overseer/pkg/materials"
	"github.com/quanxiang-cloud/overseer/pkg/reconciler"
)

var ctrLog = ctrl.Log.WithName("setup")

// Controller controller object
type Controller struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (c *Controller) SetupWithManager(mgr ctrl.Manager) error {
	if err := (&reconciler.OverseerReconciler{
		Client: c.Client,
		Scheme: c.Scheme,
		PipelineRun: &materials.PipelineRun{
			Client: c.Client,
			Scheme: c.Scheme,
		},
		Builder: &materials.Builder{
			Client: c.Client,
			Scheme: c.Scheme,
		},
		Serving: &materials.Serving{
			Client: c.Client,
			Scheme: c.Scheme,
		},
	}).SetupWithManager(mgr); err != nil {
		ctrLog.Error(err, "unable to create controller", "controller", "Overseer")
		return err
	}

	if err := (&reconciler.BuilderReconciler{
		Client: c.Client,
		Scheme: c.Scheme,
		Shipwright: &materials.Shipwright{
			Client: c.Client,
			Scheme: c.Scheme,
		},
	}).SetupWithManager(mgr); err != nil {
		ctrLog.Error(err, "unable to create controller", "controller", "Builder")
		return err
	}
	if err := (&reconciler.ServingReconciler{
		Client: c.Client,
		Scheme: c.Scheme,
		Knative: &materials.Knative{
			Client: c.Client,
			Scheme: c.Scheme,
		},
		Service: &materials.Service{
			Client: c.Client,
			Scheme: c.Scheme,
		},
	}).SetupWithManager(mgr); err != nil {
		ctrLog.Error(err, "unable to create controller", "controller", "Serving")
		return err
	}

	return nil
}
