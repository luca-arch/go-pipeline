package loop

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type (
	parallel struct {
		routine []routine
		logger  *zap.Logger
		policy  string
	}

	routine struct {
		idx      int
		task     Task
		parallel *parallel
		policy   string
	}
)

// Run executes multiple routines concurrently.
func (p parallel) Run(ctx context.Context) error {
	var wgr sync.WaitGroup

	length := len(p.routine)
	if length == 0 {
		return nil
	}

	p.logger.Info("starting", zap.Int("num", length))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Routines return results asynchronously
	resCh := make(chan error, 1)

	wgr.Add(length)

	// Spawn routines
	for _, coro := range p.routine {
		go func(task routine) {
			defer wgr.Done()

			cErr := task.Run(ctx)
			_, notify := policyCtl(cErr, task.policy)

			if notify {
				resCh <- cErr
			}
		}(coro)
	}

	// Force the return of a nil error when all routines without notifying.
	go func() {
		wgr.Wait()

		p.logger.Debug("routines completed")
		resCh <- nil
	}()

	for {
		select {
		case err := <-resCh:
			if err != nil {
				p.logger.Info("done", zap.String("result", err.Error()))
			} else {
				p.logger.Info("done")
			}

			return err
		case <-ctx.Done():
			p.logger.Info("context done")

			return nil
		}
	}
}

// AddTask inserts a new routine in the parallel queue.
func (p *parallel) AddTask(task Task, policy string) {
	p.routine = append(p.routine, routine{
		idx:      len(p.routine),
		parallel: p,
		policy:   policy,
		task:     task,
	})
}

// WithLogger sets up the logger.
func (p *parallel) WithLogger(logger *zap.Logger) *parallel {
	if logger == nil {
		logger = zap.NewNop()
	}

	p.logger = logger

	return p
}

// Run executes a single concurrent task.
func (r routine) Run(ctx context.Context) error {
	logger := r.parallel.logger.With(zap.Int("routine", r.idx))

	loop := Loop(r.task).
		WithPolicy(r.policy).
		WithLogger(logger)

	err := loop.Run(ctx)

	_, notify := policyCtl(err, r.policy)
	if notify {
		return err
	}

	if err != nil {
		logger.Info("ignoring error", zap.String("err", err.Error()))
	}

	return nil
}

// Parallel returns an executable loop.
func Parallel(tasks []Task) *parallel {
	inst := &parallel{logger: nil, policy: "", routine: nil}
	for _, task := range tasks {
		inst.AddTask(task, ExitPolicyPropagateIfErr)
	}

	return inst.WithLogger(nil)
}
