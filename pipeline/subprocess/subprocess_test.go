package subprocess_test

import (
	"context"
	"testing"

	"github.com/luca-arch/go-pipeline/pipeline/subprocess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenStreams(t *testing.T) {
	t.Parallel()

	type (
		fields struct {
			process *subprocess.Proc
		}

		want struct {
			err string
		}
	)

	testTable := map[string]struct {
		fields
		want
	}{
		"No op": {
			fields: fields{
				process: &subprocess.Proc{},
			},
			want: want{},
		},
		"With /dev/null": {
			fields: fields{
				process: &subprocess.Proc{
					Stderr: "devnul",
					Stdout: "devnul",
				},
			},
			want: want{},
		},
		"With invalid standard error": {
			fields: fields{
				process: &subprocess.Proc{
					Stderr: "/",
				},
			},
			want: want{
				err: "cannot open stderr: cannot open /: open /: is a directory",
			},
		},
		"With invalid standard output": {
			fields: fields{
				process: &subprocess.Proc{
					Stdout: "/tmp/",
				},
			},
			want: want{
				err: "cannot open stdout: cannot open /tmp/: open /tmp/: is a directory",
			},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			process := unit.fields.process
			stderr, stdout, err := process.OpenStreams()

			if unit.want.err != "" {
				assert.EqualError(t, err, unit.want.err)

				return
			}

			require.NoError(t, err)
			require.NoError(t, stderr.Close())
			require.NoError(t, stdout.Close())
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	type (
		fields struct {
			process func() *subprocess.Proc
		}

		want struct {
			err string
		}
	)

	testTable := map[string]struct {
		fields
		want
	}{
		"With single run": {
			fields: fields{
				process: func() *subprocess.Proc {
					return &subprocess.Proc{Args: []string{"-c", "exit 0"}, Command: "/bin/sh"}
				},
			},
			want: want{},
		},
		"With double run": {
			fields: fields{
				process: func() *subprocess.Proc {
					proc := &subprocess.Proc{Args: []string{"-c", "exit 0"}, Command: "/bin/sh"}
					err := proc.Run(ctx)

					require.NoError(t, err)

					return proc
				},
			},
			want: want{},
		},
		"With error code": {
			fields: fields{
				process: func() *subprocess.Proc {
					return &subprocess.Proc{Args: []string{"-c", "exit 8"}, Command: "/bin/sh"}
				},
			},
			want: want{
				err: "exit status 8",
			},
		},
		"With error": {
			fields: fields{
				process: func() *subprocess.Proc {
					return &subprocess.Proc{Command: "/bin/bon/ban"}
				},
			},
			want: want{
				err: "fork/exec /bin/bon/ban: no such file or directory",
			},
		},
		"With invalid stream": {
			fields: fields{
				process: func() *subprocess.Proc {
					return &subprocess.Proc{Command: "true", Stderr: "/not/a/file"}
				},
			},
			want: want{
				err: "stream error: cannot open stderr: cannot open /not/a/file: open /not/a/file: no such file or directory",
			},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := unit.fields.process().Run(ctx)

			if unit.want.err != "" {
				assert.EqualError(t, err, unit.want.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
