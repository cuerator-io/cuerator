package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cuerator-io/cuerator/internal/collection"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "render <collection-dir> [<input-file>]",
		Short: "Render a Cuerator collection to Kubernetes manifests",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			collectionDir := args[0]
			inputFile := ""
			if len(args) > 1 {
				inputFile = args[1]
			}
			cmd.SilenceUsage = true

			instances, err := collection.Load(collectionDir)
			if err != nil {
				return fmt.Errorf("unable to load collection: %w", err)
			}

			var inputs map[string]any

			if inputFile != "" {
				f, err := os.Open(inputFile)
				if err != nil {
					return fmt.Errorf("unable to open input file: %w", err)
				}
				defer f.Close()

				if err := json.
					NewDecoder(f).
					Decode(&inputs); err != nil {
					return fmt.Errorf("unable to decode inputs: %w", err)
				}
			}

			outputs, err := collection.Resolve(inputs, instances)

			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")

			if err := enc.Encode(outputs); err != nil {
				return fmt.Errorf("unable to encode outputs: %w", err)
			}

			return nil
		},
	}

	Root.AddCommand(cmd)
}
