package subprocess_test

import (
	"os"
	"testing"

	"github.com/luca-arch/go-pipeline/pipeline/subprocess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteCloser(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			stream []string
		}

		want struct {
			closeError string
			openError  string
			writeError string
		}
	)

	testTable := map[string]struct {
		args
		want
	}{
		"With file": {
			args: args{
				stream: tmpFile(t),
			},
		},
		"With /dev/stderr": {
			args: args{
				stream: []string{"stderr"},
			},
		},
		"With /dev/stdout": {
			args: args{
				stream: []string{"stdout"},
			},
		},
		"With /dev/null": {
			args: args{
				stream: []string{"", "devnul"},
			},
		},
		"With empty stream name": {
			args: args{
				stream: []string{"", ""},
			},
			want: want{
				openError: "missing stream name",
			},
		},
		"With non existent file": {
			args: args{
				stream: []string{"/does/not/exist"},
			},
			want: want{
				openError: "cannot open /does/not/exist: open /does/not/exist: no such file or directory",
			},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wrc, err := subprocess.WriteCloser(unit.args.stream...)
			if unit.want.openError != "" {
				assert.Nil(t, wrc)
				assert.EqualError(t, err, unit.want.openError)

				return
			}
			require.NoError(t, err)

			written, err := wrc.Write([]byte("hello"))
			if unit.want.writeError != "" {
				assert.EqualError(t, err, unit.want.writeError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 5, written)
			}

			err = wrc.Close()
			if unit.want.closeError != "" {
				assert.EqualError(t, err, unit.want.closeError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func tmpFile(t *testing.T) []string {
	t.Helper()

	fileName, err := os.CreateTemp("", "test-write-closer-*")
	require.NoError(t, err)

	return []string{fileName.Name()}
}
