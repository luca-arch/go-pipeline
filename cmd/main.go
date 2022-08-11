package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/luca-arch/go-pipeline/pipeline"
	"github.com/luca-arch/go-pipeline/template"
	"github.com/pkg/errors"
)

// errExitCode is the default exit code.
const errExitCode = 125

// task represents a task that can run.
type task interface {
	Run(context.Context) error
}

func Pipeline() (task, error) {
	args := CliArgs()

	if len(args.TemplateFiles) == 1 {
		tplCtx, err := template.Context(args.ContextFiles, args.Overrides)
		if err != nil {
			return nil, err
		}

		output, err := template.Print(args.TemplateFiles[0], tplCtx)
		if err != nil {
			return nil, err
		}

		if args.PrintOnly {
			fmt.Printf(output)
			os.Exit(0)
		}

		return pipeline.New(output)
	}

	if len(args.PipelineFiles) == 1 {
		return pipeline.NewFromFile(args.PipelineFiles[0])
	}

	return nil, errors.New("either specify a pipeline or a template file")
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
	task, err := Pipeline()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	code := Run(task)

	os.Exit(code)
}
