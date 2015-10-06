// +build dev

package main

import (
	"html/template"
	"net/http"
)

func executeTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	t, err := template.ParseFiles("templates/" + templateName + ".html.go")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}
