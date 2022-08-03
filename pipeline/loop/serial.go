package loop

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type serial struct {
	tasks  []Task
	logger *zap.Logger
	policy string
}

// Run executes multiple routines sequentially.
func (s serial) Run(ctx context.Context) error {
	s.logger.Info("starting")
	defer s.logger.Info("done")

	for _, task := range s.tasks {
		s.logger.Debug("task", zap.Any("task", task))

		err := task.Run(ctx)
		if err == nil {
			continue
		}

		_, notify := policyCtl(err, s.policy)
		if notify {
			return errors.Wrap(err, "iteration aborted")
		}

		s.logger.Warn("unreported", zap.Error(err))
	}

	return nil
}

// WithLogger sets up the logger.
func (s *serial) WithLogger(logger *zap.Logger) *serial {
	if logger == nil {
		logger = zap.NewNop()
	}

	s.logger = logger

	return s
}

// WithPolicy changes the exit policy.
func (s *serial) WithPolicy(policy string) *serial {
	if policy == "" {
		s.policy = ExitPolicyPropagateIfErr
	} else {
		s.policy = policy
	}

	return s
}

// Serial returns an executable loop.
func Serial(tasks []Task) *serial {
	inst := &serial{
		logger: nil,
		policy: "",
		tasks:  tasks,
	}

	return inst.WithLogger(nil).WithPolicy("")
}
