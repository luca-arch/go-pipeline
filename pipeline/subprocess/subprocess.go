// package subprocess wraps OS commands.
package subprocess

import (
	"context"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

// Proc represents an OS command.
type Proc struct {
	Args    []string
	Command string
	Stderr  string
	Stdout  string
}

// OpenStreams prepares the standard output and error streams.
func (p *Proc) OpenStreams() (io.WriteCloser, io.WriteCloser, error) {
	stderr, err := WriteCloser(p.Stderr, stderr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open stderr")
	}

	stdout, err := WriteCloser(p.Stdout, stdout)
	if err == nil {
		return stderr, stdout, nil
	}

	stderr.Close()

	return nil, nil, errors.Wrap(err, "cannot open stdout")
}

// Run executes the OS command.
func (p *Proc) Run(ctx context.Context) error {
	stderr, stdout, err := p.OpenStreams()
	if err != nil {
		return errors.Wrap(err, "stream error")
	}

	defer stderr.Close()
	defer stdout.Close()

	// nolint: gosec // ok
	cmd := exec.CommandContext(ctx, p.Command, p.Args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err = cmd.Run()
	if err == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err() // nolint:wrapcheck // not relevant
	default:
		return err // nolint:wrapcheck // not relevant
	}
}
