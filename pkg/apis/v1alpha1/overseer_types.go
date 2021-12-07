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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Stage the task it is performing
type Stage string

// When multiple phases exist at the same time,
// first pipelineRun stage，then build stage,last serving stage
const (
	// NoneStage init stage
	NoneStage Stage = ""

	// DoneState done State
	DoneStage Stage = "Done"

	// PipelineRunStage PipelineRun stage
	PipelineRunStage Stage = "pipelineRun"

	// BuilderStage builder state
	BuilderStage Stage = "builder"

	// ServingStage serving stage
	ServingStage Stage = "sering"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OverseerSpec defines the desired state of Overseer
type OverseerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ServiceAccountName If the deep CRD does not set serviceAccountName,
	// it will inherit from here.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	Versatile []Versatile `json:"versatile,omitempty"`

	// More info: https://kubernetes.io/docs/concepts/storage/volumes
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	Volumes []corev1.Volume `json:"volumes,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,1,rep,name=volumes"`
}

const (
	// Cancel specifies that the resource has cancel.
	Cancel = "Cancel"
)

// Type used for defining the conditiont Type field flavour
type Type string

const (
	// Succeeded specifies that the resource has finished.
	// For resources that run to completion.
	Succeeded Type = "Succeeded"
)

// Conditions condition set
type Conditions []Condition

// Condition defines a readiness condition.
// +k8s:deepcopy-gen=true
type Condition struct {
	// Type of condition
	// +required
	Type Type `json:"type" description:"type of status condition"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// We use VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic
	// differences (all other things held constant).
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" description:"last time the condition transit from one status to another"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}

//  VersatileStatus holds the versatile status
type VersatileStatus struct {
	// Ref the name of the related task.
	// +require
	Ref string `json:"ref,omitempty"`

	// StartTime is the time the task is actually started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`

	// Conditions the latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`
}

// IsFinish return ture while status is not conditionUnknown
func (v *VersatileStatus) IsFinish() bool {
	return !(v.Status == corev1.ConditionUnknown || v.Status == "")
}

var (
	// NonePhase none phase
	NonePhase = Phase{}

	// DonePhase done phase
	DonePhase = Phase{
		Stage: DoneStage,
	}
)

// Phase record the task being performed
type Phase struct {
	Phase string `json:"phase,omitempty"`
	Stage Stage  `json:"stage,omitempty"`
}

func (p *Phase) IsNone() bool {
	return p.Phase == "" && p.Stage == NoneStage
}

func (p *Phase) IsDone() bool {
	return p.Stage == DoneStage
}

func NewPhase(label string, stage Stage) Phase {
	return Phase{
		Phase: label,
		Stage: stage,
	}
}

// OverseerStatus defines the observed state of Overseer
type OverseerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase Phase `json:"phase,omitempty"`

	// StartTime is the time the task is actually started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the task completed.
	// +optional
	CompletionTime metav1.Time `json:"completionTime,omitempty"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// +optional
	VersatileStatus []VersatileStatus `json:"versatileStatus,omitempty"`
}

func (o *OverseerStatus) SetVersatileStatus(vs VersatileStatus) {
	if o.VersatileStatus == nil {
		o.VersatileStatus = []VersatileStatus{vs}
		return
	}
	last := len(o.VersatileStatus) - 1
	if o.VersatileStatus[last].Ref == vs.Ref {
		o.VersatileStatus[last] = vs
		return
	}

	o.VersatileStatus = append(o.VersatileStatus, vs)
}

// SetFalse If any one task fails, the overall failure.
// this overseer will failure and stop.
func (o *OverseerStatus) SetFalse() {
	o.Status = corev1.ConditionFalse
	o.CompletionTime = metav1.NewTime(time.Now())
}

// SetSuccess all task done.
func (o *OverseerStatus) SetSuccess() {
	o.Status = corev1.ConditionTrue
	o.Phase = DonePhase
	o.CompletionTime = metav1.NewTime(time.Now())
}

func (o *OverseerStatus) IsUnkonwn() bool {
	return o.Status == corev1.ConditionUnknown ||
		o.Status == ""
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase.phase`
//+kubebuilder:printcolumn:name="Stage",type=string,JSONPath=`.status.phase.stage`
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Overseer is the Schema for the overseers API
type Overseer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OverseerSpec   `json:"spec,omitempty"`
	Status OverseerStatus `json:"status,omitempty"`
}

func (o *Overseer) getNextPhase() Phase {
	if o.Status.Phase.IsDone() {
		return o.Status.Phase
	}

	var (
		index int
		stage Stage
	)

	phase := o.Status.Phase

	if !o.Status.Phase.IsNone() {
		for i, versatile := range o.Spec.Versatile {
			if versatile.Name == phase.Phase {
				index = i
			}
		}
	}

	stage = phase.Stage
	for ; index < len(o.Spec.Versatile); index++ {
		versatile := &o.Spec.Versatile[index]
		// the current task has not yet executed stage
		switch stage {
		case NoneStage:
			if versatile.PipelineRun != nil {
				stage = PipelineRunStage
				break
			}
			fallthrough
		case PipelineRunStage:
			if versatile.Builder != nil {
				stage = BuilderStage
				break
			}
			fallthrough
		case BuilderStage:
			if versatile.Serving != nil {
				stage = ServingStage
				break
			}
			fallthrough
		case ServingStage:
			stage = NoneStage
			continue
		}
		return NewPhase(versatile.Name, stage)
	}

	return DonePhase
}

func (o *Overseer) GetVersatile() *Versatile {
	phase := o.Status.Phase

	if phase.IsNone() {
		phase = o.getNextPhase()
	}

	o.Status.Phase = phase
	for _, versatile := range o.Spec.Versatile {
		if versatile.Name == phase.Phase {
			return &versatile
		}
	}

	return nil
}

// IsDone if status is unkonwn where return false
func (o *Overseer) IsDone() bool {
	return !o.Status.IsUnkonwn()
}

// ShoudContinue return to true if the current task is not completed,
// or there are still unexecuted tasks.
func (o *Overseer) ShoudContinue() bool {
	if len(o.Spec.Versatile) == 0 {
		return false
	}

	if o.Status.Phase.IsDone() {
		return false
	}

	if len(o.Status.VersatileStatus) == 0 ||
		o.Status.VersatileStatus[len(o.Status.VersatileStatus)-1].IsFinish() {
		//get the index of the next task
		nextPhase := o.getNextPhase()
		if nextPhase == DonePhase {
			return false
		}
		o.Status.Phase = nextPhase
	}
	return true
}

//+kubebuilder:object:root=true

// OverseerList contains a list of Overseer
type OverseerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Overseer `json:"items"`
}
