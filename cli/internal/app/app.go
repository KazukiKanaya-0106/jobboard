package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kanaya/jobboard-cli/internal/config"
	"github.com/kanaya/jobboard-cli/internal/hub"
	"github.com/kanaya/jobboard-cli/internal/runner"
	"github.com/kanaya/jobboard-cli/internal/slack"
)

const (
	statusCompleted = "completed"
	statusFailed    = "failed"
)

type App struct {
	config *config.Config
	hub    *hub.Client
	slack  *slack.Notifier
	runner *runner.Runner
}

func New(config *config.Config, hub *hub.Client, slack *slack.Notifier, runner *runner.Runner) *App {
	return &App{
		config: config,
		hub:    hub,
		slack:  slack,
		runner: runner,
	}
}

func (app *App) Run(ctx context.Context) int {
	startedAt := time.Now()
	var hubStarted bool

	if app.config.Hub.Enabled() {
		if err := app.hub.Start(ctx, startedAt); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to notify Hub start: %v\n", err)
		} else {
			hubStarted = true
		}
	}

	result, runErr := app.runner.Run(ctx, app.config.Execution.Command)
	if runErr != nil && result == nil {
		fmt.Fprintf(os.Stderr, "error: failed to execute command: %v\n", runErr)
		result = &runner.Result{ExitCode: 1, Error: runErr}
	}
	if result == nil {
		result = &runner.Result{ExitCode: 1}
	}

	finishedAt := time.Now()
	status := statusCompleted
	if result.ExitCode != 0 || result.Error != nil {
		status = statusFailed
	}

	if hubStarted {
		if err := app.hub.Finish(ctx, status, finishedAt, finishedAt.Sub(startedAt)); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to notify Hub finish: %v\n", err)
		}
	}

	if app.config.Slack.Enabled() {
		var errorText string
		if status == statusFailed {
			switch {
			case result.Error != nil:
				errorText = result.Error.Error()
			case result.Stderr != "":
				errorText = result.Stderr
			default:
				errorText = fmt.Sprintf("Exit code: %d", result.ExitCode)
			}
		}

		payload := slack.Payload{
			Command:    strings.Join(app.config.Execution.Command, " "),
			Tag:        app.config.Hub.Tag,
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
			Duration:   finishedAt.Sub(startedAt),
			Status:     status,
			ExitCode:   result.ExitCode,
			Error:      errorText,
		}

		if err := app.slack.Notify(ctx, payload); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to send Slack notification: %v\n", err)
		}
	}

	return result.ExitCode
}
