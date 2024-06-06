package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"
	"net/http"
	"strconv"
	"time"
)

func CalendarHandler(w http.ResponseWriter, r *http.Request) {

	// Comprobamos si el usuario esta logeado
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si el año es menor del 24 o superior al 27 dejamos el 24
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2024 || year > 2027 {
		year = time.Now().Year()
	}

	weeks := models.GetWeeksOfYear(year)
	utils.Tmpl.ExecuteTemplate(w, "calendar.html", struct {
		Weeks        []models.Week
		SelectedYear int
	}{
		Weeks:        weeks,
		SelectedYear: year,
	})

}
