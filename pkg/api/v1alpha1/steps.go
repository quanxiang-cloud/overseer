package v1alpha1

// StepSpec define task execution steps.
type StepSpec struct {
	// Name is the name of this steo.
	Name string `json:"name,omitempty"`

	// Description is a user-facing description of the step.
	// +optional
	Description string `json:"description,omitempty"`

	// Template declares yaml file of kubernetes resource
	Template string `json:"template,omitempty"`
}
