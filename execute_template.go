// +build !dev

package main

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.html.go"))
)

func executeTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	templates.ExecuteTemplate(w, templateName+".html.go", nil)
}
