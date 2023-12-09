package installationcrd

import (
	"github.com/cuerator-io/cuerator/crd/internal/status"
)

// Status is the status portion of an [Installation] resource.
type Status struct {
	Conditions status.ConditionSet `json:"conditions,omitempty"`
	Tag        Tag                 `json:"tag,omitempty"`
}

// Tag describes the image tag that an [Installation]'s version constraint
// resolves to.
type Tag struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Digest  string `json:"digest,omitempty"`
}
