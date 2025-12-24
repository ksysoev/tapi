package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksysoev/tapi/pkg/cmd"
)

var (
	version = "dev"
	name    = "tapi"
)

func main() {
	os.Exit(runApp())
}

func runApp() int {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	command := cmd.InitCommand(cmd.BuildInfo{
		Version: version,
		AppName: name,
	})

	if err := command.ExecuteContext(ctx); err != nil {
		return 1
	}

	return 0
}
