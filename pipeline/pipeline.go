// package pipeline provides a way to load a procedural list of tasks from a YAML file.
package pipeline

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// NewFromFile reads a YAML file and parse it as a Node.
func NewFromFile(file string) (*Node, error) {
	str, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open %s", file)
	}

	exec, err := New(string(str))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid YAML in %s", file)
	}

	return exec, nil
}

// New parses a YAML string and returns the root node.
func New(str string) (*Node, error) {
	var exec Node

	if err := yaml.Unmarshal([]byte(str), &exec); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal")
	}

	logger, err := exec.LogConfig.Logger()
	if err != nil {
		return nil, errors.Wrap(err, "cannot create logger")
	}

	return exec.WithLogger(logger), nil
}
