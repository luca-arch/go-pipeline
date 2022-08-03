package subprocess

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	// devnul is a string alias to `/dev/null`.
	devnul = "devnul"
	// mode is how files are open.
	mode = os.O_APPEND | os.O_CREATE | os.O_RDWR // nolint:nosnakecase // go package
	// perm is new files' permissions.
	perm = 0o644
	// stderr is a string alias to `/dev/stderr`.
	stderr = "stderr"
	// stdout is a string alias to `/dev/stdout`.
	stdout = "stdout"
)

// writeCloser implements io.WriteCloser.
type writeCloser struct {
	c io.WriteCloser
	w io.Writer
}

// Write writes into the stream.
func (s writeCloser) Write(p []byte) (int, error) {
	// nolint:wrapcheck // not relevant
	return s.w.Write(p)
}

// Close closes the stream.
func (s writeCloser) Close() error {
	if s.c == nil {
		return nil
	}

	// nolint:wrapcheck // not relevant
	return s.c.Close()
}

// WriteCloser returns a new open stream.
func WriteCloser(args ...string) (io.WriteCloser, error) {
	for _, stream := range args {
		switch stream {
		case "":
			continue
		case devnul:
			return writeCloser{nil, io.Discard}, nil
		case stderr:
			return writeCloser{nil, os.Stderr}, nil
		case stdout:
			return writeCloser{nil, os.Stdout}, nil
		default:
			f, err := os.OpenFile(stream, mode, perm)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot open %s", stream)
			}

			return writeCloser{f, f}, nil
		}
	}

	return nil, errors.New("missing stream name")
}
