package main

import (
	"html/template"
	"log"
    "strconv"
	"net/http"
    "math"
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

        split := split(total, 3)

        log.Println("Split:", split)

		templates.ExecuteTemplate(w, "result.html", total)
	} else {
		log.Println("Split endpoint called without post command")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func toNearestCent(val float64) float64 {
    return math.Round(val*100) / 100
}

func split(amount float64, participants int) []float64 {
    if participants <= 0 {
        return []float64{amount}
    }

    baseAmount := toNearestCent(amount / float64(participants))
    splits := make([]float64, participants)

    for i := range splits {
        splits[i] = baseAmount
    }

    totalAssigned := baseAmount * float64(participants)

    remainingAmount := toNearestCent(amount - totalAssigned)

    for i := 0; remainingAmount > 0 && i < participants; i++ {
        splits[i] = toNearestCent(splits[i] + 0.01)
        remainingAmount = toNearestCent(remainingAmount - 0.01)
    }

    return splits
}
