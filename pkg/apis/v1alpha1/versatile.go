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

package v1alpha1

import corev1 "k8s.io/api/core/v1"

type Versatile struct {
	// Name If the name is empty, will get oveseer name.
	// +optional
	Name string `json:"name,omitempty"`

	// PipelineRun Relying on tekton's pipeline to complete various tasks
	// +optional
	PipelineRun *PipelineRunSpec `json:"pipelineRun,omitempty"`

	Builder *BuilderSpec `json:"builder,omitempty"`

	Serving *ServingSpec `json:"serving,omitempty"`
}

type PipelineRunSpec struct {
	PipelineRef string          `json:"pipelineRef,omitempty"`
	Params      []Param         `json:"params,omitempty"`
	Workspace   []Workspace     `json:"workspace,omitempty"`
	Template    *corev1.PodSpec `json:"template,omitempty"`
}

type Workspace struct {
	// Name is the name of the workspace populated by the volume.
	Name string `json:"name"`

	// SubPath is optionally a directory on the volume which should be used
	// for this binding (i.e. the volume will be mounted at this sub directory).
	// +optional
	SubPath string `json:"subPath,omitempty"`
}
