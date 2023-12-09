package installationcrd

// Spec is the specification section of an [Installation].
type Spec struct {
	Image             string      `json:"image"`
	VersionConstraint string      `json:"versionConstraint"`
	Inputs            []InputSpec `json:"inputs"`
}

// InputSpec describes a single input to an [Installation]'s CUE templates.
type InputSpec struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
