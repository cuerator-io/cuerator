package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a Cuerator collection to Kubernetes manifests",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx := cmd.Context()
			// dir := args[0]
			cmd.SilenceUsage = true

			return nil
		},
	}

	Root.AddCommand(cmd)
}
