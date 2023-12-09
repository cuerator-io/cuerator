package reconciler

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/cuerator-io/cuerator/crd/installationcrd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

// version describes an image tag that can be treated as a semantic version.
type version struct {
	TagName string
	SemVer  *semver.Version
}

func (r *Reconciler) resolveVersion(
	ctx context.Context,
	inst *installationcrd.Installation,
) (err error) {
	cond := metav1.Condition{Type: "VersionResolved"}
	defer r.setConditionOnReturn(ctx, inst, &cond, &err)

	cond.Reason = "InvalidImageReference"

	ref, err := registry.ParseReference(inst.Spec.Image)
	if err != nil {
		return err
	}

	if ref.Reference != "" {
		return errors.New("image name must not include a tag or digest")
	}

	cond.Reason = "InvalidVersionConstraint"

	con, err := semver.NewConstraint(inst.Spec.VersionConstraint)
	if err != nil {
		return err
	}

	cond.Reason = "RegistryError"

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

	cond.Reason = "ConstraintNotSatisifed"

	if len(tags) == 0 {
		return fmt.Errorf(
			"the %q repository has no tags",
			ref.Repository,
		)
	}

	versions := parseTags(tags)
	if len(versions) == 0 {
		return fmt.Errorf(
			"the %q repository has %d tags, none of which can be parsed as semantic versions",
			ref.Repository,
			len(tags),
		)
	}

	ver, ok := selectLatestVersion(versions, con)
	if !ok {
		return fmt.Errorf(
			"the %q repository has %d tags, %d of which can be parsed as semantic versions, but none match the %q constraint",
			ref.Repository,
			len(tags),
			len(versions),
			con,
		)
	}

	cond.Reason = "RegistryError"

	d, err := repo.Resolve(ctx, ver.TagName)
	if err != nil {
		return err
	}

	cond.Reason = "ConstraintSatisfied"

	inst.Status.Tag = installationcrd.TagStatus{
		Name:    ver.TagName,
		Version: ver.SemVer.String(),
		Digest:  d.Digest.String(),
	}

	return nil
}

// parseTags parses the given tags as semantic versions, omitting any that
// are invalid.
func parseTags(tags []string) []version {
	var result []version

	for _, tag := range tags {
		ver, err := semver.NewVersion(tag)
		if err == nil {
			result = append(result, version{tag, ver})
		}
	}

	return result
}

// selectLatestVersion returns the latest (highest) of the given versions that
// satisfies the given constraint.
func selectLatestVersion(
	candidates []version,
	constraint *semver.Constraints,
) (version, bool) {
	var matches []version
	for _, v := range candidates {
		if constraint.Check(v.SemVer) {
			matches = append(matches, v)
		}
	}

	slices.SortFunc(
		matches,
		func(a, b version) int {
			return b.SemVer.Compare(a.SemVer)
		},
	)

	for _, tag := range matches {
		return tag, true
	}

	return version{}, false
}
