package models

import (
	"CalendarApp/utils"
	"strconv"
	"time"
)

type Week struct {
	StartDate   string
	EndDate     string
	WeekID      string
	ButtonClass string
	Year        int
}

func GetWeeksOfYear(year int) []Week {

	// Variables a utilizar
	var weeks []Week
	t := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	weekNumber := 1
	now := time.Now()

	// Para cada dia de la semana añadimos un numero
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, 1)
	}

	// Para todos los años cuando es su año
	for t.Year() == year {

		startOfWeek := t
		endOfWeek := t.AddDate(0, 0, 6)
		buttonClass := "btn-primary" // Color del boton
		if endOfWeek.Before(now) {
			buttonClass = "btn-secondary"
		}

		// Establecemos el ID de la semana
		weekID := strconv.Itoa(year) + strconv.Itoa(weekNumber)

		// Lista de datos de las semanas
		weeks = append(weeks, Week{
			StartDate:   utils.FormatDate(startOfWeek),
			EndDate:     utils.FormatDate(endOfWeek),
			WeekID:      weekID,
			ButtonClass: buttonClass,
			Year:        year,
		})
		t = t.AddDate(0, 0, 7)
		weekNumber++
	}
	return weeks

}
