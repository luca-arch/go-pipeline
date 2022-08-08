package template

import (
	"os"
	"text/template"
)

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"isDir": func(path string) bool {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return false
			}

			return fileInfo.IsDir()
		},
		"isFile": func(path string) bool {
			fileInfo, err := os.Stat(path)
			if err != nil {
				return false
			}

			return !fileInfo.IsDir()
		},
	}
}
