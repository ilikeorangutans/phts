package web

import (
	"html/template"
	"log"
	"path/filepath"
)

func GetTemplates(files ...string) *template.Template {
	name := filepath.Base(files[0])
	result, err := template.New(name).ParseFiles(files...)
	if err != nil {
		log.Println(err)
	}

	return result
}
