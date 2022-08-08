package template

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

func Context(files []string, vars []string) (map[string]interface{}, error) {
	context := make(map[string]interface{})

	for _, ctxFile := range files {
		tmpCtx := make(map[string]interface{})

		yml, err := os.ReadFile(ctxFile)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(yml, &tmpCtx); err != nil {
			return nil, err
		}

		if err := mergo.Merge(&context, tmpCtx); err != nil {
			return nil, err
		}
	}

	for n := range vars {
		s := strings.SplitN(vars[n], "=", 2) // nolint:gomnd // legit
		context[s[0]] = s[1]
	}

	return context, nil
}

func Print(templateFile string, context map[string]interface{}) (string, error) {
	templateText, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateFile).
		Funcs(templateFunctions()).
		Parse(string(templateText))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, context)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
