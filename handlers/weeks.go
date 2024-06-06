package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"
	"net/http"
	"strconv"
	"time"
)

func WeekWorkersHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica si el usuario está logueado
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Obtiene el ID de la semana y el año de los parámetros de la URL
	weekID := r.URL.Query().Get("weekID")
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" || weekID == "" {
		http.Error(w, "Missing year or weekID", http.StatusBadRequest)
		return
	}

	// Convierte el año a un entero y maneja posibles errores
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2024 || year > 2027 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	// Convierte el ID de la semana a un entero y maneja posibles errores
	weekNumber, err := strconv.Atoi(weekID)
	if err != nil {
		http.Error(w, "Invalid week number in week ID", http.StatusBadRequest)
		return
	}

	// Obtiene la lista de trabajadores desde la base de datos
	workers, err := models.GetWorkers()

	if err != nil {
		http.Error(w, "Unable to retrieve workers", http.StatusInternalServerError)
		return
	}

	// Calcula las fechas de inicio y fin de la semana basada en el número de semana y año
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	startOfWeek := startOfYear.AddDate(0, 0, (weekNumber-1)*7)
	for startOfWeek.Weekday() != time.Monday {
		startOfWeek = startOfWeek.AddDate(0, 0, -1)
	}
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	// Define los días de la semana
	days := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
	var weekSchedule []models.DaySchedule

	// Crea el horario semanal
	for i, day := range days {
		date := startOfWeek.AddDate(0, 0, i)
		weekSchedule = append(weekSchedule, models.DaySchedule{
			Day:     day,
			Date:    utils.FormatDate(date),
			Workers: workers,
		})
	}

	// Estructura de datos para pasar a la plantilla
	tmplData := struct {
		WeekID       string
		Year         string
		StartDate    string
		EndDate      string
		WeekSchedule []models.DaySchedule
	}{
		WeekID:       weekID,
		Year:         yearStr,
		StartDate:    utils.FormatDate(startOfWeek),
		EndDate:      utils.FormatDate(endOfWeek),
		WeekSchedule: weekSchedule,
	}

	// Renderiza la plantilla de la semana
	utils.Tmpl.ExecuteTemplate(w, "weeks.html", tmplData)
}
