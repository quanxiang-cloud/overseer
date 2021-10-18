package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Condition defines a readiness condition
type Condition struct {
	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`

	// +required
	ResourceRef map[string]StepCondition `json:"resourceRef,omitempty"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

func (c *Condition) IsNil() bool {
	return c.Status == ""
}

func (c *Condition) Init() {
	c.Status = corev1.ConditionUnknown
	c.ResourceRef = make(map[string]StepCondition)
}

func (c *Condition) IsUnknown() bool {
	return c.Status == corev1.ConditionUnknown
}

func (c *Condition) IsTrue() bool {
	return c.Status == corev1.ConditionTrue
}

func (c *Condition) IsFalse() bool {
	return c.Status == corev1.ConditionFalse
}

func (c *Condition) IsFinish(name string) bool {
	if c.ResourceRef == nil {
		return true
	}

	ref, ok := c.ResourceRef[name]
	if !ok {
		return true
	}

	return ref.IsFinish()
}

type StepConditionType string

const (
	StepConditionUnknown = "Unknown"
	StepConditionPending = "Pending"
	StepConditionRunning = "Running"
	StepConditionSuccess = "Success"
	StepConditionFail    = "Fail"
)

type StepCondition struct {
	GroupVersionKind string `json:"groupVersionKind,omitempty"`

	RefName string `json:"refName,omitempty"`

	State StepConditionType `json:"state,omitempty"`

	// A human-readable message indicating details about why the volume is in this state.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// Reason is a brief string that describes any failure and is meant
	// for machine parsing and tidy display in the CLI.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
}

func (s StepCondition) IsFinish() bool {
	return s.State == StepConditionSuccess ||
		s.State == StepConditionFail
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
