package main

import (
	"html/template"
	"log"
    "strconv"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func splitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		totalStr := r.FormValue("total")
		total, err := strconv.ParseFloat(totalStr, 64)
		if err != nil {
			http.Error(w, "Invalid total amount", http.StatusBadRequest)
			return
		}

		templates.ExecuteTemplate(w, "result.html", total)
	} else {
		log.Println("Split endpoint called without post command")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
