package loop_test

import (
	"context"
	"errors"
)

var (
	errA = errors.New("error A")
	errB = errors.New("error B")
)

type testTask struct {
	err error
}

func (t testTask) Run(_ context.Context) error {
	return t.err
}
