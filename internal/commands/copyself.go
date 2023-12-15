package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "copy-self",
		Short: "Create a copy of the Cuerator executable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			to := args[0]
			cmd.SilenceUsage = true

			from, err := os.Executable()
			if err != nil {
				return fmt.Errorf("unable to determine executable path: %w", err)
			}

			r, err := os.Open(from)
			if err != nil {
				return fmt.Errorf("unable to open source file: %w", err)
			}
			defer r.Close()

			if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
				return fmt.Errorf("unable to create destination directory: %w", err)
			}

			w, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("unable to create destination file: %w", err)
			}
			defer w.Close()

			if _, err := io.Copy(w, r); err != nil {
				return fmt.Errorf("unable to copy file contents: %w", err)
			}

			return nil
		},
	}

	Root.AddCommand(cmd)
}
