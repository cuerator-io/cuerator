package commands

import (
	"context"

	"github.com/cuerator-io/cuerator/internal/di"
	"github.com/cuerator-io/cuerator/internal/operator/installationcrd"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func init() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Cuerator Kubernetes operator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cmd.SilenceUsage = true

			return imbue.Invoke2(
				ctx,
				di.Container,
				func(
					ctx context.Context,
					m manager.Manager,
					r *installationcrd.Reconciler,
				) error {
					err := builder.
						ControllerManagedBy(m).
						For(&installationcrd.Installation{}).
						WithEventFilter(predicate.GenerationChangedPredicate{}).
						Complete(r)

					if err != nil {
						return err
					}

					return m.Start(ctx)
				},
			)
		},
	}

	Root.AddCommand(cmd)
}
