package installationcrd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/semver"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconciler reconciles the state of [Installation] resources.
type Reconciler struct {
	Manager manager.Manager
	Client  client.Client
	Logger  logr.Logger
}

// version describes an image tag that can be treated as a semantic version.
type version struct {
	TagName string
	SemVer  *semver.Version
}

// Reconcile performs a full reconciliation for the object referred to by the
// request.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	inst := &Installation{}
	if err := r.Client.Get(ctx, req.NamespacedName, inst); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.resolveVersion(ctx, inst); err != nil {
		return reconcile.Result{}, err
	}

	// if err := r.reconcileJob(ctx, inst); err != nil {
	// 	return reconcile.Result{}, err
	// }

	return reconcile.Result{
		RequeueAfter: 5 * time.Minute,
	}, nil
}

func (r *Reconciler) setConditionOnReturn(
	ctx context.Context,
	inst *Installation,
	cond *metav1.Condition,
	err *error,
) {
	if cond.Status != metav1.ConditionUnknown && cond.Reason == "" {
		panic("the condition must have a reason")
	}

	if *err == nil {
		cond.Status = metav1.ConditionTrue
	} else {
		cond.Message = (*err).Error()
		cond.Status = metav1.ConditionFalse
	}

	cond.ObservedGeneration = inst.Generation
	inst.Status.Conditions.Merge(*cond)

	*err = errors.Join(
		*err,
		r.Client.Status().Update(ctx, inst),
	)
}

func (r *Reconciler) resolveVersion(
	ctx context.Context,
	inst *Installation,
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
			inst.Spec.VersionConstraint,
		)
	}

	cond.Reason = "RegistryError"

	d, err := repo.Resolve(ctx, ver.TagName)
	if err != nil {
		return err
	}

	cond.Reason = "ConstraintSatisfied"

	inst.Status.Tag = Tag{
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
