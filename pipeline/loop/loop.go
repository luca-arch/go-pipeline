// package loop provides some methods to run a Task
package loop

import (
	"context"

	"go.uber.org/zap"
)

const (
	ExitPolicyNone           = "none"
	ExitPolicyRestart        = "restart"
	ExitPolicyRestartIfErr   = "restart-if-err"
	ExitPolicyPropagate      = "propagate"
	ExitPolicyPropagateIfErr = "propagate-if-err"
)

type (
	loop struct {
		task   Task
		logger *zap.Logger
		policy string
	}

	Task interface {
		Run(ctx context.Context) error
	}
)

// Run executes the loop (restarts with ExitPolicyRestart/ExitPolicyRestartIfErr).
func (l loop) Run(ctx context.Context) error {
	l.logger.Debug("starting loop", zap.String("policy", l.policy))
	defer l.logger.Debug("closing loop", zap.String("policy", l.policy))

	for {
		err := l.task.Run(ctx)
		restart, notify := policyCtl(err, l.policy)

		switch {
		case notify:
			return err // nolint:wrapcheck // legit
		case restart && err != nil:
			l.logger.Info("ignoring error", zap.String("err", err.Error()))
		case !restart:
			return nil
		}
	}
}

// WithLogger sets up the logger.
func (l *loop) WithLogger(logger *zap.Logger) *loop {
	if logger == nil {
		logger = zap.NewNop()
	}

	l.logger = logger

	return l
}

// WithPolicy changes the exit policy.
func (l *loop) WithPolicy(policy string) *loop {
	if policy == "" {
		l.policy = ExitPolicyPropagateIfErr
	} else {
		l.policy = policy
	}

	return l
}

// Loop constructor.
func Loop(task Task) *loop {
	inst := &loop{
		logger: nil,
		policy: "",
		task:   task,
	}

	return inst.WithLogger(nil).WithPolicy("")
}

// policyCtl returns whether to restart and/or notify a result.
func policyCtl(err error, policy string) (bool, bool) {
	switch policy {
	case ExitPolicyNone:
		return false, false
	case ExitPolicyRestart:
		return true, false
	case ExitPolicyRestartIfErr:
		return err != nil, false
	case ExitPolicyPropagate:
		return false, true
	default: // ExitPolicyPropagateIfErr
		return false, err != nil
	}
}
