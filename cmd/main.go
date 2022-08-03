package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"bitbucket.org/lucacontini/z6/pipeline"
	"github.com/pkg/errors"
)

// errExitCode is the default exit code.
const errExitCode = 125

// task represents a task that can run.
type task interface {
	Run(context.Context) error
}

// Run executes a task and returns any propagated exit code.
func Run(e task) int {
	var exErr *exec.ExitError

	err := e.Run(context.Background())

	switch {
	case err == nil:
		return 0
	case errors.As(err, &exErr):
		return exErr.ExitCode()
	default:
		log.Printf("FATAL: %v", err)
	}

	return errExitCode
}

func main() {
	task, err := pipeline.NewFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("error %v", err)
	}

	code := Run(task)

	os.Exit(code)
}
