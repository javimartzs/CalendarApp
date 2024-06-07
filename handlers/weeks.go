package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const dataDir = "data"

func init() {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.Mkdir(dataDir, 0755)
		if err != nil {
			log.Fatalf("Error creating data directory: %v", err)
		}
	}
}

func saveTableDataToFile(weekID, year string, data map[string]interface{}) error {
	fileName := filepath.Join(dataDir, weekID+"_"+year+".json")
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, dataJSON, 0644)
}

func loadTableDataFromFile(weekID, year string) (map[string]interface{}, error) {
	fileName := filepath.Join(dataDir, weekID+"_"+year+".json")
	dataJSON, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func deleteFutureTableFiles(currentDate time.Time) error {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			parts := strings.Split(file.Name(), "_")
			if len(parts) >= 2 {
				weekID := parts[0]
				weekStartDate, err := time.Parse("2006-01-02", weekID)
				if err != nil {
					continue
				}

				if weekStartDate.After(currentDate) {
					err = os.Remove(filepath.Join(dataDir, file.Name()))
					if err != nil {
						log.Printf("Error deleting file %s: %v", file.Name(), err)
					}
				}
			}
		}
	}
	return nil
}

func SaveTableStateHandler(w http.ResponseWriter, r *http.Request) {
	var tableData struct {
		WeekID  string                 `json:"weekID"`
		Year    string                 `json:"year"`
		Data    map[string]interface{} `json:"data"`
		Summary map[string]float64     `json:"summary"`
	}
	if err := json.NewDecoder(r.Body).Decode(&tableData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := saveTableDataToFile(tableData.WeekID, tableData.Year, map[string]interface{}{
		"data":    tableData.Data,
		"summary": tableData.Summary,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoadTableStateHandler(w http.ResponseWriter, r *http.Request) {
	weekID := r.URL.Query().Get("weekID")
	year := r.URL.Query().Get("year")

	data, err := loadTableDataFromFile(weekID, year)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "No data found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// New endpoint to handle worker changes
func UpdateWorkersHandler(w http.ResponseWriter, r *http.Request) {
	var workerData struct {
		WeekID  string          `json:"weekID"`
		Year    string          `json:"year"`
		Workers []models.Worker `json:"workers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&workerData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := loadTableDataFromFile(workerData.WeekID, workerData.Year)
	if err != nil {
		if os.IsNotExist(err) {
			data = make(map[string]interface{})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data["workers"] = workerData.Workers

	err = saveTableDataToFile(workerData.WeekID, workerData.Year, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete future table files
	currentDate := time.Now()
	err = deleteFutureTableFiles(currentDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ResetTableStateHandler(w http.ResponseWriter, r *http.Request) {
	weekID := r.URL.Query().Get("weekID")
	year := r.URL.Query().Get("year")

	fileName := filepath.Join(dataDir, weekID+"_"+year+".json")
	err := os.Remove(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
