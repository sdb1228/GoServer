// +build !dev

package main

import (
	"html/template"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.html.go"))
)

func lookupTemplate(name string) *template.Template {
	return templates.Lookup(name + ".html.go")
}
