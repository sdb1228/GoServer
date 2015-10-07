// +build dev

package main

import (
	"html/template"
)

func lookupTemplate(name string) *template.Template {
	t, _ := template.ParseFiles("templates/" + name + ".html.go")
	return t
}
