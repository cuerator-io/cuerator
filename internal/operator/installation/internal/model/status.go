package model

import "github.com/cuerator-io/cuerator/internal/operator/internal/crd"

// Status is the status portion of an [Installation] resource.
type Status struct {
	Conditions     crd.ConditionSet `json:"conditions,omitempty"`
	DesiredVersion *Version         `json:"desiredVersion,omitempty"`
}

// Version describes the image tag that an [Installation]'s version constraint
// resolves to.
type Version struct {
	Normalized string `json:"normalized,omitempty"`
	Image      string `json:"image,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Digest     string `json:"digest,omitempty"`
}

const (
	// VersionResolvedCondition is a condition type that indicates whether an
	// installations version constraint has been resolved to a specific image
	// tag.
	VersionResolvedCondition = "VersionResolved"

	// VersionResolvedReasonInvalidImage indicates that a version has not been
	// resolved because the image name is invalid.
	VersionResolvedReasonInvalidImage = "InvalidImage"

	// VersionResolvedReasonInvalidConstraint indicates that a version has not
	// been resolved because the version constraint is invalid.
	VersionResolvedReasonInvalidConstraint = "InvalidConstraint"

	// VersionResolvedReasonConstraintNotSatisfied indicates that a version has
	// not been resolved because none of the image's tags match the version
	// constraint.
	VersionResolvedReasonConstraintNotSatisfied = "ConstraintNotSatisfied"
)
