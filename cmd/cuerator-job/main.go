package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("not enough arguments")
	}

	command := os.Args[1]
	path := os.Args[2]

	switch command {
	case "init":
		return initCommand(path)
	case "render":
		return renderCommand(path)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func renderCommand(path string) error {
	fmt.Println("rendering installation templates")
	return nil
}
