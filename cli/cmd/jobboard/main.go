package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanaya/jobboard-cli/internal/app"
	"github.com/kanaya/jobboard-cli/internal/config"
	"github.com/kanaya/jobboard-cli/internal/hub"
	"github.com/kanaya/jobboard-cli/internal/runner"
	"github.com/kanaya/jobboard-cli/internal/slack"
)

func main() {
	config, warnings, err := config.Load(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	application := app.New(
		config,
		hub.NewClient(config.Hub, &http.Client{Timeout: config.Hub.Timeout}),
		slack.NewNotifier(config.Slack, &http.Client{Timeout: config.Slack.Timeout}),
		runner.New(),
	)

	exitCode := application.Run(ctx)
	os.Exit(exitCode)
}
