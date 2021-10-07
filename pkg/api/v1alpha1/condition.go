package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// Condition defines a readiness condition
type Condition struct {
	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`

	// +required
	ResourceRef map[string]StepCondition `json:"resourceRef,omitempty"`
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

func (c *Condition) IsUnFalse() bool {
	return c.Status == corev1.ConditionFalse
}

type StepConditionType string

const (
	StepConditionUnknown = "Unknown"
	StepConditionPending = "Pending"
	StepConditionRunning = "Running"
	StepConditionSuceess = "Success"
	StepConditionFail    = "Fail"
)

type StepCondition struct {
	GroupVersionKind string            `json:"groupVersionKind,omitempty"`
	State            StepConditionType `json:"state,omitempty"`
}
