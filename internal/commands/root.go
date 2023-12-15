package commands

import (
	"os"

	"github.com/cuerator-io/cuerator/internal/di"
	"github.com/dogmatiq/imbue"
	"github.com/spf13/cobra"
)

// Root is the root "bsctl" command.
var Root = &cobra.Command{
	Use:   "cuerator",
	Short: "Continuous delivery, configuration and policy enforcement for Kubernetes using CUE",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Add the currently-executing Cobra CLI command the the DI container.
		//
		// This hook is called after the CLI arguments are resolved to a
		// specific command.
		//
		// This allows other DI declarations to make use of the flags passed to
		// the command.
		imbue.With0(
			di.Container,
			func(imbue.Context) (*cobra.Command, error) {
				return cmd, nil
			},
		)
	},
}

func init() {
	Root.SetIn(os.Stdin)
	Root.SetOut(os.Stdout)
	Root.SetErr(os.Stderr)
}
