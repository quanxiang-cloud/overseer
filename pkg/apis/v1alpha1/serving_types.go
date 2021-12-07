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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Knative struct {
}

type Service struct {
	Ports []corev1.ServicePort `json:"ports,omitempty"`
}

type Serve struct {
	// Knative
	// +optional
	Knative *Knative `json:"knative,omitempty"`

	// Service
	// +optional
	Service *Service `json:"service,omitempty"`
}

// ServingSpec defines the desired state of Serving
type ServingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name string `json:"name,omitempty"`

	// ServiceAccountName
	ServiceAccountName *string `json:"serviceAccountName,omitempty"`

	// Serve Serve
	Serve `json:",inline"`

	Template *corev1.PodSpec `json:"template,omitempty"`
}

// ServingStatus defines the observed state of Serving
type ServingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Ref string `json:"ref,omitempty"`

	// StartTime is the time the task is actually started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the task completed.
	// +optional
	CompletionTime metav1.Time `json:"completionTime,omitempty"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// Conditions the latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

func (s *ServingStatus) IsKnown() bool {
	return s.Status == "" || s.Status == corev1.ConditionUnknown
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Serving is the Schema for the servings API
type Serving struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServingSpec   `json:"spec,omitempty"`
	Status ServingStatus `json:"status,omitempty"`
}

func (s *Serving) IsDone() bool {
	return !s.Status.IsKnown()
}

func (s *Serving) GetCondition(_t Type) *Condition {
	for _, cond := range s.Status.Conditions {
		if cond.Type == _t {
			return cond.DeepCopy()
		}
	}
	return nil
}

//+kubebuilder:object:root=true

// ServingList contains a list of Serving
type ServingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Serving `json:"items"`
}
