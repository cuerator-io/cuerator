package installation

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

// version describes an image tag that can be treated as a semantic version.
type version struct {
	Tag     string
	Version *semver.Version
	Digest  string
}

func (i *installer) resolveVersion(ctx context.Context) error {
	ref, err := registry.ParseReference(i.Installation.Spec.Image)
	if err != nil {
		i.Installation.Status.Conditions.Set(
			metav1.Condition{
				Type:    model.VersionResolvedCondition,
				Status:  metav1.ConditionFalse,
				Reason:  model.VersionResolvedReasonInvalidImage,
				Message: err.Error(),
			},
		)
		return nil
	}

	if ref.Reference != "" {
		i.Installation.Status.Conditions.Set(
			metav1.Condition{
				Type:    model.VersionResolvedCondition,
				Status:  metav1.ConditionFalse,
				Reason:  model.VersionResolvedReasonInvalidImage,
				Message: "image name must not include a tag or digest",
			},
		)
		return nil
	}

	con, err := semver.NewConstraint(i.Installation.Spec.VersionConstraint)
	if err != nil {
		i.Installation.Status.Conditions.Set(
			metav1.Condition{
				Type:    model.VersionResolvedCondition,
				Status:  metav1.ConditionFalse,
				Reason:  model.VersionResolvedReasonInvalidConstraint,
				Message: err.Error(),
			},
		)
		return nil
	}

	reg, err := remote.NewRegistry(ref.Host())
	if err != nil {
		return err
	}

	repo, err := reg.Repository(ctx, ref.Repository)
	if err != nil {
		return err
	}

	tags, err := registry.Tags(ctx, repo)
	if err != nil {
		return err
	}

	ver, ok := highestMatchingVersion(con, tags)
	if !ok {
		i.Installation.Status.Conditions.Set(
			metav1.Condition{
				Type:   model.VersionResolvedCondition,
				Status: metav1.ConditionFalse,
				Reason: model.VersionResolvedReasonConstraintNotSatisfied,
				Message: fmt.Sprintf(
					"the %q image does not have any tags that match the %q constraint",
					i.Installation.Spec.Image,
					i.Installation.Spec.VersionConstraint,
				),
			},
		)
		return nil
	}

	desc, err := repo.Resolve(ctx, ver.Tag)
	if err != nil {
		return err
	}

	i.Installation.Status.Conditions.Set(
		metav1.Condition{
			Type:   model.VersionResolvedCondition,
			Status: metav1.ConditionTrue,
			Message: fmt.Sprintf(
				"the %q image tag %q matches the %q constraint",
				ver.Tag,
				i.Installation.Spec.Image,
				i.Installation.Spec.VersionConstraint,
			),
		},
	)

	i.Installation.Status.DesiredVersion = &model.Version{
		Image:      i.Installation.Spec.Image,
		Tag:        ver.Tag,
		Digest:     desc.Digest.String(),
		Normalized: ver.Version.String(),
	}

	return nil
}

// highestMatchingVersion returns the "highest" (latest) version that matches the given
// constraint, or false if no versions match.
//
// Any tag names that are not valid semantic versions are ignored.
func highestMatchingVersion(
	con *semver.Constraints,
	tags []string,
) (version, bool) {
	var latest version

	for _, tag := range tags {
		ver, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}

		if !con.Check(ver) {
			continue
		}

		if latest.Version == nil || ver.GreaterThan(latest.Version) {
			latest.Tag = tag
			latest.Version = ver
		}
	}

	return latest, latest.Version != nil
}
