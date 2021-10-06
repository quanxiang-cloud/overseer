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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OverseerSpec defines the desired state of Overseer
type OverseerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Description is a user-facing description of the oversser.
	// +optional
	Description string `json:"description,omitempty"`

	// Params declares a list of input parameters that must be supplied when
	// this overseer is run.
	// +optional
	Params []ParamSpec `json:"params,omitempty"`

	// Steps are the steps of the overseer.
	Steps []StepSpec `json:"steps,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Overseer is the Schema for the overseers API
type Overseer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OverseerSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// OverseerList contains a list of Overseer
type OverseerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Overseer `json:"items"`
}
