package pipeline

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// LogConfig.
type LogConfig struct {
	Debug    bool   `yaml:"debug"`
	Disabled bool   `yaml:"disabled"`
	Level    string `yaml:"level"`

	inst *zap.Logger
}

// Logger returns a zap Logger instance.
func (l *LogConfig) Logger() (*zap.Logger, error) {
	if l.inst != nil {
		return l.inst, nil
	}

	var (
		err    error
		logger *zap.Logger
	)

	switch {
	case l.Disabled:
		return zap.NewNop(), nil
	case l.Debug:
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, errors.Wrap(err, "zap")
	}

	l.inst = logger

	return logger, nil
}
