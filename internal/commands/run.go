package commands

import (
	"context"

	"github.com/cuerator-io/cuerator/internal/di"
	"github.com/cuerator-io/cuerator/internal/operator/installation"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Cuerator Kubernetes operator",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cmd.SilenceUsage = true

			return imbue.Invoke1(
				ctx,
				di.Container,
				func(
					ctx context.Context,
					m manager.Manager,
				) error {
					if err := installation.BuildController(m); err != nil {
						return err
					}

					return m.Start(ctx)
				},
			)
		},
	}

	Root.AddCommand(cmd)
}
