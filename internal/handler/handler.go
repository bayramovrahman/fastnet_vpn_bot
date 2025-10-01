package handler

import (
	"html/template"
	"net/http"
)

// Home handler
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/templates/home.page.tmpl")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, nil)
}
