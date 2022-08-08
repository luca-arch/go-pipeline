package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/luca-arch/go-pipeline/pipeline"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	type (
		fields struct {
			instance pipeline.Node
		}

		want struct {
			err error
		}
	)

	testTable := map[string]struct {
		fields
		want
	}{
		"Success": {
			fields: fields{
				instance: pipeline.Node{
					Command: "true",
				},
			},
			want: want{},
		},
		"Errored": {
			fields: fields{
				instance: pipeline.Node{
					Command: "false",
					Name:    "failure-1",
				},
			},
			want: want{
				err: errors.New("task failure-1: exit status 1"),
			},
		},
		"With children": {
			fields: fields{
				instance: pipeline.Node{
					Steps: []pipeline.Node{
						{
							Command: "true",
						},
						{
							Command: "echo",
						},
					},
				},
			},
			want: want{},
		},
		"With parallel and exit policies": {
			fields: fields{
				instance: pipeline.Node{
					Parallel: []pipeline.Node{
						{
							Command: "sh",
							Args:    []string{"-c", "sleep 0.1 && exit 1"},
							OnExit:  "restart",
						},
						{
							Command: "sh",
							Args:    []string{"-c", "sleep 0.5 && exit 4"},
							OnExit:  "report-if-err",
						},
					},
				},
			},
			want: want{
				err: errors.New("task parallel: task sh: exit status 4"),
			},
		},
		"With file": {
			fields: fields{
				instance: load(t, "../testdata/test-pipeline-001.yaml"),
			},
			want: want{},
		},
		"With file (error)": {
			fields: fields{
				instance: load(t, "../testdata/test-pipeline-002.yaml"),
			},
			want: want{
				err: errors.New("task test-pipeline-002: iteration aborted: task parallel: task sh: exit status 67"),
			},
		},
		"With file (timeout)": {
			fields: fields{
				instance: load(t, "../testdata/test-pipeline-004.yaml"),
			},
			want: want{
				err: errors.New("task test-pipeline-004: iteration aborted: task exit-64: context deadline exceeded"),
			},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			err := unit.fields.instance.WithLogger(nil).Run(ctx)

			if unit.want.err != nil {
				assert.EqualError(t, err, unit.want.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func load(t *testing.T, file string) pipeline.Node {
	t.Helper()

	n, err := pipeline.NewFromFile(file)
	require.NoError(t, err)

	n = n.WithLogger(nil)

	return *n
}
