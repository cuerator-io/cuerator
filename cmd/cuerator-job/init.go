package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func initCommand(path string) error {
	return copySelf(filepath.Join(path, "bin"))
}

func copySelf(dir string) error {
	from, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to determine executable path: %w", err)
	}

	to := filepath.Join(dir, filepath.Base(from))

	r, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("unable to open source file: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return fmt.Errorf("unable to create destination directory: %w", err)
	}

	w, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}
	defer w.Close()

	if err := os.Chmod(to, 0755); err != nil {
		return fmt.Errorf("unable to set destination file permissions: %w", err)
	}

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("unable to copy file contents: %w", err)
	}

	return nil
}
