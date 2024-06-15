package handlers

import (
	"ajaxbits.com/bsplit/internal/splits"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func RootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func SplitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		totalStr := r.FormValue("total")
		total, err := strconv.ParseFloat(totalStr, 64)
		if err != nil {
			http.Error(w, "Invalid total amount", http.StatusBadRequest)
			return
		}

		split := splits.Split(total, 3, splits.Even)

		log.Println("Split:", split)

		templates.ExecuteTemplate(w, "result.html", total)
	} else {
		log.Println("Split endpoint called without post command")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
