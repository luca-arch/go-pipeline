// package pipeline provides a way to load a procedural list of tasks from a YAML file.
package pipeline

import (
	"context"
	"time"

	"bitbucket.org/lucacontini/z6/pipeline/loop"
	"bitbucket.org/lucacontini/z6/pipeline/subprocess"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Node represents the pipeline execution.
type Node struct {
	Args      []string      `yaml:"args,flow"`
	Command   string        `yaml:"path"`
	LogConfig LogConfig     `yaml:"log"`
	Name      string        `yaml:"name"`
	OnExit    string        `yaml:"onExit"`
	Parallel  []Node        `yaml:"parallel,flow"`
	Stderr    string        `yaml:"stderr"`
	Stdout    string        `yaml:"stdout"`
	Steps     []Node        `yaml:"steps,flow"`
	Timeout   time.Duration `yaml:"timeout"`

	logger *zap.Logger
}

// ID returns the identifier (name) of the current node.
func (n *Node) ID() string {
	switch {
	case n.Name != "":
		return n.Name
	case n.IsCommand():
		return n.Command
	case n.IsParallel():
		return "parallel"
	case n.IsSerial():
		return "serial"
	default:
		return "[empty]"
	}
}

// IsCommand returns whether the node represents an executable command.
func (n *Node) IsCommand() bool {
	return n.Command != ""
}

// IsParallel returns whether the node represents a list of parallel tasks.
func (n *Node) IsParallel() bool {
	return len(n.Parallel) > 0
}

// IsCommand returns whether the node represents a list of tasks.
func (n *Node) IsSerial() bool {
	return len(n.Steps) > 0
}

// Run executes the pipeline.
func (n Node) Run(ctx context.Context) error {
	ctl := ctx

	if n.Timeout > 0 {
		n.logger.Info("set timeout", zap.Any("seconds", n.Timeout))
		c, cancel := context.WithTimeout(ctx, n.Timeout)
		ctl = c

		defer cancel()
	}

	if err := n.Task().Run(ctl); err != nil {
		return errors.Wrapf(err, "task %s", n.ID())
	}

	return nil
}

// Task return the current node as a loop.
func (n *Node) Task() loop.Task { //nolint:ireturn // Legit interface
	n.logBogusConfig()
	n.propagateLogger()

	switch {
	case n.IsCommand():
		cmd := &subprocess.Proc{
			Args:    n.Args,
			Command: n.Command,
			Stderr:  n.Stderr,
			Stdout:  n.Stdout,
		}

		return loop.Loop(cmd).WithLogger(n.logger).WithPolicy(n.OnExit)
	case n.IsParallel():
		tasks := typecast(n.Parallel)

		// Maybe TODO? n.OnExit has no effect on this node.
		return loop.Parallel(tasks).WithLogger(n.logger)
	case !n.IsSerial():
		n.logger.Warn("noop node")
	}

	tasks := typecast(n.Steps)

	return loop.Serial(tasks).WithLogger(n.logger).WithPolicy(n.OnExit)
}

// WithLogger sets up the logger.
func (n *Node) WithLogger(logger *zap.Logger) *Node {
	if logger == nil {
		logger = zap.NewNop()
	}

	n.logger = logger.With(zap.String("task", n.ID()))

	return n
}

// logBogusConfig warns about any bogus node configuration.
func (n *Node) logBogusConfig() {
	switch {
	case n.IsCommand() && n.IsParallel():
		n.logger.Warn("bogus `parallel` list")
	case n.IsCommand() && n.IsSerial():
		n.logger.Warn("bogus `steps` list")
	case n.IsParallel() && n.IsSerial():
		n.logger.Warn("bogus `steps` list")
	}
}

// propagateLogger attaches the node logger instance to its children.
func (n *Node) propagateLogger() {
	for i := range n.Parallel {
		n.Parallel[i].WithLogger(n.logger)
	}

	for i := range n.Steps {
		n.Steps[i].WithLogger(n.logger)
	}
}

// typecast converts a list of nodes into a loop.
func typecast(nodes []Node) []loop.Task {
	loopTasks := make([]loop.Task, 0)
	for _, task := range nodes {
		loopTasks = append(loopTasks, task)
	}

	return loopTasks
}
