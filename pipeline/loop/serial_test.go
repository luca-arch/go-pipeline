package loop_test

import (
	"context"
	"testing"
	"time"

	"github.com/luca-arch/go-pipeline/pipeline/loop"
	"github.com/stretchr/testify/assert"
)

func TestSerialRun(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			policy string
			tasks  []loop.Task
		}

		want struct {
			err string
		}
	)

	testTable := map[string]struct {
		args
		want
	}{
		"Success": {
			args: args{
				tasks: []loop.Task{
					testTask{nil},
					testTask{nil},
				},
			},
			want: want{},
		},
		"Empty": {
			args: args{
				tasks: []loop.Task{},
			},
			want: want{},
		},
		"With error": {
			args: args{
				tasks: []loop.Task{
					testTask{nil},
					testTask{errA},
					testTask{nil},
				},
			},
			want: want{
				err: "iteration aborted: error A",
			},
		},
		"With exit policy (none)": {
			args: args{
				policy: loop.ExitPolicyNone,
				tasks: []loop.Task{
					testTask{nil},
					testTask{errB},
					testTask{nil},
					testTask{errB},
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

			err := loop.Serial(unit.args.tasks).WithPolicy(unit.args.policy).Run(ctx)

			if unit.want.err != "" {
				assert.EqualError(t, err, unit.want.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
