package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cuerator-io/cuerator/internal/commands"
	"github.com/dogmatiq/ferrite"
)

var version = "0.0.0"

func main() {
	ferrite.Init()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	commands.Root.Version = version

	if err := commands.Root.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
