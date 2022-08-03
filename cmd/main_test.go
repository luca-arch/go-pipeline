package main_test

import (
	"testing"

	main "bitbucket.org/lucacontini/z6/cmd"
	"bitbucket.org/lucacontini/z6/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctional(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			file string
		}

		want struct {
			code int
		}
	)

	testTable := map[string]struct {
		args
		want
	}{
		"File test-pipeline-001.yaml": {
			args: args{
				file: "../testdata/test-pipeline-001.yaml",
			},
			want: want{code: 0},
		},
		"File test-pipeline-002.yaml": {
			args: args{
				file: "../testdata/test-pipeline-002.yaml",
			},
			want: want{code: 67},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			task := mkPipeline(t, unit.args.file)
			code := main.Run(task)

			assert.Equal(t, unit.want.code, code)
		})
	}
}

func mkPipeline(t *testing.T, file string) *pipeline.Node {
	t.Helper()

	p, err := pipeline.NewFromFile(file)
	require.NoError(t, err)

	return p.WithLogger(nil)
}
