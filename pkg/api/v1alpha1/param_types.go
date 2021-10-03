package v1alpha1

// ParamSpec defines parameters needed beyond typed inputs.
type ParamSpec struct {

	// Name declares the name by which a parameter is referenced.
	Name string `json:"name,omitempty"`

	// Description is a user-facing description of the parameter.
	// +optional
	Description string `json:"description,omitempty"`

	// Default is the value a parameter takes if no input value is supplied.
	// +optional
	Default string `json:"default,omitempty"`
}

// Param declares an value to use for the parameter called name.
type Param struct {
	Name string `json:"name,omitempty"`

	// Value is from default or user-input.
	Value string `json:"value,omitempty"`
}
