package installation

import (
	"github.com/cuerator-io/cuerator/internal/operator"
	"github.com/cuerator-io/cuerator/internal/operator/installation/internal/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// BuildController builds a controller that manages [model.Installation]
// resources.
func BuildController(
	m manager.Manager,
) error {
	b := &scheme.Builder{
		GroupVersion: schema.GroupVersion{
			Group:   operator.GroupName,
			Version: model.APIVersion,
		},
	}

	b.Register(
		&model.Installation{},
		&model.InstallationList{},
	)

	if err := b.AddToScheme(m.GetScheme()); err != nil {
		return err
	}

	c := &reconciler{
		Manager: m,
		Logger:  m.GetLogger().V(1).WithName("installation"),
	}

	return builder.
		ControllerManagedBy(m).
		For(&model.Installation{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(c)
}
