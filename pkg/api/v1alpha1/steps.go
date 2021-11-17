package v1alpha1

// StepSpec define task execution steps.
type StepSpec struct {
	// Name is the name of this step.
	Name string `json:"name,omitempty"`

	// Separate if separate is true,the resource is no longer under jurisdiction
	Separate bool `json:"separate,omitempty"`

	// Description is a user-facing description of the step.
	// +optional
	Description string `json:"description,omitempty"`

	// Template declares yaml file of kubernetes resource
	Template string `json:"template,omitempty"`
}
