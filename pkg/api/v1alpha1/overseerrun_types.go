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

// OverseerRunSpec defines the desired state of OverseerRun
type OverseerRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	OverseerRef *OverseerRef `json:"overseerRef,omitempty"`

	// Params is a list of parameter names and values.
	Params []Param `json:"params,omitempty"`
}

// OverseerRunStatus defines the observed state of OverseerRun
type OverseerRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Status `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.condition.status`
//+kubebuilder:printcolumn:name="Messge",type=string,JSONPath=`.status.condition.message`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// OverseerRun is the Schema for the overseerruns API
type OverseerRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OverseerRunSpec   `json:"spec,omitempty"`
	Status OverseerRunStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OverseerRunList contains a list of OverseerRun
type OverseerRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OverseerRun `json:"items"`
}

// OverseerRef can be used to refer to a specific instance of a Overseer
type OverseerRef struct {
	// Name of the referent
	Name string `json:"name,omitempty"`
}

type Status struct {
	Condition Condition `json:"condition,omitempty" `
}
