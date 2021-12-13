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

// Shipwright more info https://github.com/shipwright-io/build
type Shipwright struct {

	// Strategy references the BuildStrategy to use to build the container image.
	Strategy ShipwrightStrategy `json:"strategy,omitempty"`
}

//  ShipwrightStrategy can be used to refer to a specific instance of a buildstrategy.
type ShipwrightStrategy struct {
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name"`

	// BuildStrategyKind indicates the kind of the buildstrategy, namespaced or cluster scoped.
	Kind string `json:"kind,omitempty"`
}

// BuilderEngine various builder collections.
type BuilderEngine struct {
	Shipwright *Shipwright `json:"shipwright,omitempty"`
}

// BuilderSpec defines the desired state of Builder
type BuilderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ServiceAccountName
	ServiceAccountName *string `json:"serviceAccountName,omitempty"`

	// Git
	Git Git `json:"git,omitempty"`

	// Image
	Image Image `json:"image,omitempty"`

	// Param
	Params []Param `json:"params,omitempty"`

	// BuilderEngine
	BuilderEngine `json:",omitempty"`
}

// BuilderStatus defines the observed state of Builder
type BuilderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ref the name of builder.
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

func (b *BuilderStatus) IsKnown() bool {
	return b.Status == "" || b.Status == corev1.ConditionUnknown
}

//+genclient
//+genclient:noStatus
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ref",type=string,JSONPath=`.status.ref`
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Builder is the Schema for the builders API
type Builder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuilderSpec   `json:"spec,omitempty"`
	Status BuilderStatus `json:"status,omitempty"`
}

func (b *Builder) IsDone() bool {
	return !b.Status.IsKnown()
}

func (b *Builder) GetCondition(_t Type) *Condition {
	for _, cond := range b.Status.Conditions {
		if cond.Type == _t {
			return cond.DeepCopy()
		}
	}
	return nil
}

//+kubebuilder:object:root=true

// BuilderList contains a list of Builder
type BuilderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Builder `json:"items"`
}
