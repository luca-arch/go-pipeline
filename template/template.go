package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type cliContextOverride []string

func (c *cliContextOverride) Set(value string) error {
	*c = append(*c, value)

	return nil
}

func (c *cliContextOverride) String() string {
	return ""
}

func Parse(templateFile, contextFile string, overrides cliContextOverride) {
	funcMap := template.FuncMap{
		// "title": strings.Title,
	}

	context := make(map[string]interface{})
	if contextFile != "" {
		if err := yaml.Unmarshal([]byte(mustOpenFile(contextFile)), &context); err != nil {
			log.Fatalf("%s: %s", contextFile, err)
		}
	}

	for n := range overrides {
		s := strings.SplitN(overrides[n], "=", 2) // nolint:gomnd // legit
		context[s[0]] = s[1]
	}

	templateText := mustOpenFile(templateFile)
	if templateText == "" {
		log.Fatalf("%s: empty template", templateFile)
	}

	tmpl, err := template.New(templateText).Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	err = tmpl.Execute(os.Stdout, context)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}

func main() {
	var overrides cliContextOverride

	flag.Var(&overrides, "set", "Context overrides")

	templateFile := flag.String("template-file", "", "Path to template file")
	contextFile := flag.String("context-file", "", "Path to context file")

	flag.Parse()

	Parse(*templateFile, *contextFile, overrides)
}

func mustOpenFile(file string) string {
	body, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return string(body)
}
