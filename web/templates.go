package web

import (
	"html/template"
	"log"
	"path/filepath"

	humanize "github.com/dustin/go-humanize"
)

var templateFunctions = template.FuncMap{
	"humanizeTime": humanize.Time,
}

func GetTemplates(files ...string) *template.Template {
	name := filepath.Base(files[0])
	result, err := template.New(name).Funcs(templateFunctions).ParseFiles(files...)
	if err != nil {
		log.Println(err)
	}

	return result
}
