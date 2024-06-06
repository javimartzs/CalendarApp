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

	// Obtiene la lista de semanas del año
	weeks := models.GetWeeksOfYear(year)

	// Encuentra la semana correspondiente al weekID
	var selectedWeek models.Week
	for _, week := range weeks {
		if week.WeekID == weekID {
			selectedWeek = week
			break
		}
	}

	if selectedWeek.WeekID == "" {
		http.Error(w, "Week not found", http.StatusNotFound)
		return
	}

	// Obtiene la lista de trabajadores desde la base de datos
	workers, err := models.GetWorkers()
	if err != nil {
		http.Error(w, "Unable to retrieve workers", http.StatusInternalServerError)
		return
	}

	// Define los días de la semana
	days := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
	var weekSchedule []models.DaySchedule

	// Convierte las fechas de string a time.Time
	startOfWeek, _ := time.Parse("2006-01-02", selectedWeek.StartDate)

	// Crea el horario semanal
	for i, day := range days {
		date := startOfWeek.AddDate(0, 0, i)
		weekSchedule = append(weekSchedule, models.DaySchedule{
			Day:     day,
			Date:    utils.FormatDate(date), // Asegúrate de que utils.FormatDate formatea correctamente la fecha
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
		StartDate:    selectedWeek.StartDate,
		EndDate:      selectedWeek.EndDate,
		WeekSchedule: weekSchedule,
	}

	// Renderiza la plantilla de la semana
	utils.Tmpl.ExecuteTemplate(w, "weeks.html", tmplData)
}
