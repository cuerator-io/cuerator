package installationcrd

// Spec is the specification section of an [Installation].
type Spec struct {
	Image             string             `json:"image"`
	VersionConstraint string             `json:"versionConstraint"`
	Inputs            []Input            `json:"inputs"`
	InputsFrom        []InputsFromSource `json:"inputsFrom"`
}

// Input is a single input to an [Installation].
type Input struct {
	Path      string       `json:"path"`
	Value     any          `json:"value,omitempty"`
	ValueFrom *InputSource `json:"valueFrom,omitempty"`
}

// InputSource describes how to obtain the value of an input from a specific key
// within a ConfigMap or Secret.
type InputSource struct {
	ConfigMapKeyRef *InputKeySelector `json:"configMapKeyRef,omitempty"`
	SecretKeyRef    *InputKeySelector `json:"secretKeyRef,omitempty"`
}

// InputKeySelector describes a key within a ConfigMap or Secret.
type InputKeySelector struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Optional bool   `json:"optional,omitempty"`
	Static   bool   `json:"static,omitempty"`
}

// InputsFromSource describes how to obtain the values of inputs from a
// ConfigMap or Secret.
type InputsFromSource struct {
	ConfigMapRef *InputRef `json:"configMapRef,omitempty"`
	SecretRef    *InputRef `json:"secretRef,omitempty"`
	Path         string    `json:"path,omitempty"`
	Static       bool      `json:"static,omitempty"`
}

// InputRef is a reference to a ConfigMap or Secret.
type InputRef struct {
	Name     string `json:"name,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}
