package installationcrd

import "github.com/cuerator-io/cuerator/internal/operator/internal/crd"

// Status is the status portion of an [Installation] resource.
type Status struct {
	Conditions crd.ConditionSet `json:"conditions,omitempty"`
	Tag        *Tag             `json:"tag,omitempty"`
}

// Tag describes the image tag that an [Installation]'s version constraint
// resolves to.
type Tag struct {
	Image             string `json:"image,omitempty"`
	Name              string `json:"name,omitempty"`
	Digest            string `json:"digest,omitempty"`
	NormalizedVersion string `json:"normalizedVersion,omitempty"`
}
