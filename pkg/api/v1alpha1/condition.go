package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *OverseerRunStatus) IsNil() bool {
	return o.Status == ""
}

func (o *OverseerRunStatus) Init() {
	o.Status = corev1.ConditionUnknown
	o.ResourceRef = make(map[string]RefCondition)
}

func (o *OverseerRunStatus) IsUnknown() bool {
	return o.Status == corev1.ConditionUnknown
}

func (o *OverseerRunStatus) IsTrue() bool {
	return o.Status == corev1.ConditionTrue
}

func (o *OverseerRunStatus) IsFalse() bool {
	return o.Status == corev1.ConditionFalse
}

func (o *OverseerRunStatus) IsFinish(name string) bool {
	if o.ResourceRef == nil {
		return true
	}

	ref, ok := o.ResourceRef[name]
	if !ok {
		return true
	}

	return ref.IsFinish()
}

type RefCondition struct {
	GroupVersionKind string `json:"groupVersionKind,omitempty"`

	Name string `json:"name,omitempty"`

	Conditions []Condition `json:"conditions,omitempty"`
}

func (r *RefCondition) IsFinish() bool {
	if len(r.Conditions) == 0 {
		return false
	}

	for _, elem := range r.Conditions {
		if elem.Status == corev1.ConditionUnknown {
			return false
		}
	}
	return true
}

func (r *RefCondition) IsFalse() bool {
	for _, elem := range r.Conditions {
		if elem.Status == corev1.ConditionFalse {
			return true
		}
	}

	return false
}

type Condition struct {
	Status corev1.ConditionStatus `json:"status,omitempty"`

	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`

	// A human-readable message indicating details about why the volume is in this state.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// Reason is a brief string that describes any failure and is meant
	// for machine parsing and tidy display in the CLI.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
}

type Phase string

const (
	PhaseDone Phase = "Done"
)

func (p *Phase) IsDone() bool {
	return *p == PhaseDone
}

func (p *Phase) Equal(val string) bool {
	return string(*p) == val
}

func (p *Phase) IsNil() bool {
	return *p == ""
}

func (p *Phase) Sting() string {
	return string(*p)
}
