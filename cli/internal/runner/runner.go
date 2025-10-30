package runner

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
)

type Runner struct{}

type Result struct {
	ExitCode int
	Stderr   string
	Error    error
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Run(ctx context.Context, command []string) (*Result, error) {
	if len(command) == 0 {
		return nil, errors.New("no command provided to run")
	}

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	var stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()

	result := &Result{
		ExitCode: 0,
		Stderr:   stderrBuf.String(),
	}

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
			result.Error = err
			return result, nil
		}
		result.ExitCode = 1
		result.Error = err
		return result, nil
	}

	return result, nil
}
