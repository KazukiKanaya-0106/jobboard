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

func (app *App) Run(ctx context.Context) (exitCode int) {
	loc := app.config.Time.Location
	if loc == nil {
		loc = time.Local
	}

	startedRaw := time.Now()
	startedAt := startedRaw.In(loc)
	var (
		hubStarted bool
		result     *runner.Result
		status     = statusCompleted
		errorText  string
	)

	defer func() {
		finishedRaw := time.Now()
		finishedAt := finishedRaw.In(loc)
		duration := finishedRaw.Sub(startedRaw)

		trimmedError := strings.TrimSpace(errorText)
		var hubErrorText *string
		if trimmedError != "" {
			hubErrorText = &trimmedError
		}

		if hubStarted {
			finishCtx := context.Background()
			if timeout := app.config.Hub.Timeout; timeout > 0 {
				var cancel context.CancelFunc
				finishCtx, cancel = context.WithTimeout(finishCtx, timeout)
				defer cancel()
			}
			if err := app.hub.Finish(finishCtx, status, finishedAt, duration, hubErrorText); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to notify Hub finish: %v\n", err)
			}
		}

		if app.config.Slack.Enabled() {
			slackCtx := context.Background()
			if timeout := app.config.Slack.Timeout; timeout > 0 {
				var cancel context.CancelFunc
				slackCtx, cancel = context.WithTimeout(slackCtx, timeout)
				defer cancel()
			}

			payload := slack.Payload{
				Command:    strings.Join(app.config.Execution.Command, " "),
				Tag:        app.config.Hub.Tag,
				StartedAt:  startedAt,
				FinishedAt: finishedAt,
				Duration:   duration,
				Status:     status,
				ExitCode:   exitCode,
				Error:      trimmedError,
			}

			if err := app.slack.Notify(slackCtx, payload); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to send Slack notification: %v\n", err)
			}
		}

		if status == statusFailed && exitCode == 0 {
			exitCode = 1
		}
	}()

	if app.config.Hub.Enabled() {
		if err := app.hub.Start(ctx, startedAt); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to notify Hub start: %v\n", err)
		} else {
			hubStarted = true
		}
	}

	res, runErr := app.runner.Run(ctx, app.config.Execution.Command)
	if runErr != nil && res == nil {
		fmt.Fprintf(os.Stderr, "error: failed to execute command: %v\n", runErr)
		res = &runner.Result{ExitCode: 1, Error: runErr}
	}
	if res == nil {
		res = &runner.Result{ExitCode: 1}
	}
	result = res

	exitCode = result.ExitCode
	if result.ExitCode != 0 || result.Error != nil || ctx.Err() != nil {
		status = statusFailed
		switch {
		case result.Error != nil:
			errorText = result.Error.Error()
		case result.Stderr != "":
			errorText = result.Stderr
		}
	}

	if ctx.Err() != nil {
		if errorText == "" {
			errorText = fmt.Sprintf("terminated by signal: %v", ctx.Err())
		}
	} else if status == statusFailed && errorText == "" {
		errorText = fmt.Sprintf("Exit code: %d", result.ExitCode)
	}

	return exitCode
}
