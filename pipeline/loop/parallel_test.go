package loop_test

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/lucacontini/z6/pipeline/loop"
	"github.com/stretchr/testify/assert"
)

func TestParallelRun(t *testing.T) {
	t.Parallel()

	type (
		fields struct {
			instance func(t *testing.T) loop.Task
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
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					return mkParallel(t, testTask{nil}, testTask{nil})
				},
			},
			want: want{},
		},
		"Empty": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					return mkParallel(t)
				},
			},
			want: want{},
		},
		"With error": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					return mkParallel(t, testTask{nil}, testTask{errA})
				},
			},
			want: want{
				err: errA,
			},
		},
		"With exit policy - none": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					p := loop.Parallel(nil)
					p.AddTask(testTask{nil}, loop.ExitPolicyNone)
					p.AddTask(testTask{errA}, loop.ExitPolicyNone)

					return p
				},
			},
			want: want{},
		},
		"With exit policy - propagate": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					p := loop.Parallel(nil)
					p.AddTask(testTask{nil}, loop.ExitPolicyPropagate)
					p.AddTask(testTask{errA}, loop.ExitPolicyNone)

					return p
				},
			},
			want: want{},
		},
		"With exit policy - propagate-if-err": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					p := loop.Parallel(nil)
					p.AddTask(testTask{nil}, loop.ExitPolicyPropagateIfErr)
					p.AddTask(testTask{errA}, loop.ExitPolicyPropagateIfErr)
					p.AddTask(testTask{errB}, loop.ExitPolicyNone)

					return p
				},
			},
			want: want{
				err: errA,
			},
		},
		"With exit policy - restart": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					p := loop.Parallel(nil)
					p.AddTask(testTask{nil}, loop.ExitPolicyRestart)
					p.AddTask(testTask{errB}, loop.ExitPolicyRestart)

					return p
				},
			},
			want: want{},
		},
		"With exit policy - restart-if-err": {
			fields: fields{
				instance: func(t *testing.T) loop.Task {
					t.Helper()

					p := loop.Parallel(nil)
					p.AddTask(testTask{errA}, loop.ExitPolicyRestartIfErr)
					p.AddTask(testTask{errB}, loop.ExitPolicyRestartIfErr)

					return p
				},
			},
			want: want{},
		},
	}

	for name, unit := range testTable {
		unit := unit

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			err := unit.fields.instance(t).Run(ctx)

			if unit.want.err != nil {
				assert.EqualError(t, err, unit.want.err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// nolint:ireturn // ok
func mkParallel(t *testing.T, tasks ...loop.Task) loop.Task {
	t.Helper()

	return loop.Parallel(tasks)
}
