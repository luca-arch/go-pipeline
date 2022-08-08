package main

import "flag"

type (
	// cliArgMulti implements flag.Var to read a cli parameter multiple times
	cliArgMulti []string

	// cliArgs cli parameter to change the application's behaviour
	cliArgs struct {
		ContextFiles  []string
		Overrides     []string
		PipelineFiles []string
		PrintOnly     bool
		TemplateFiles []string
	}
)

func (c *cliArgMulti) Set(value string) error {
	*c = append(*c, value)

	return nil
}

func (c *cliArgMulti) String() string {
	return ""
}

func CliArgs() *cliArgs {
	var (
		contextFiles  cliArgMulti
		overrides     cliArgMulti
		printOnly     bool
		templateFiles cliArgMulti
	)

	flag.Var(&contextFiles, "context", "Path to context file (can be passed multiple times)")
	flag.Var(&overrides, "set", "Override global context values (can be passed multiple times)")
	flag.Var(&templateFiles, "template", "Path to template file (can be passed multiple times)")
	flag.BoolVar(&printOnly, "print-only", false, "Only parse the template and print its output, do not execute")

	flag.Parse()

	return &cliArgs{
		ContextFiles:  contextFiles,
		Overrides:     overrides,
		PipelineFiles: flag.Args(),
		PrintOnly:     printOnly,
		TemplateFiles: templateFiles,
	}
}
