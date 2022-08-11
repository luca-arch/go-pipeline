package template

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type tplContext map[string]interface{}

func Context(files []string, vars []string) (tplContext, error) {
	context := make(tplContext)

	for _, ctxFile := range files {
		tmpCtx := make(tplContext)

		yml, err := os.ReadFile(ctxFile)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(yml, &tmpCtx); err != nil {
			return nil, err
		}

		context = mergeMaps(context, tmpCtx)
	}

	for n := range vars {
		s := strings.SplitN(vars[n], "=", 2) // nolint:gomnd // legit
		context[s[0]] = s[1]
	}

	return context, nil
}

func Print(templateFile string, context tplContext) (string, error) {
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

func mergeMaps(b, a tplContext) tplContext {
	out := make(tplContext, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(tplContext); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(tplContext); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
