package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"

	"net/http"
	"strconv"
)

func WorkerProfileHandler(w http.ResponseWriter, r *http.Request) {

	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Path[len("/worker/"):])
	if err != nil {
		http.Error(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		if err := models.DeleteWorker(id); err != nil {
			http.Error(w, "Unable to delete worker", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	worker, err := models.GetWorkerByID(id)
	if err != nil {
		http.Error(w, "Unable to retrieve worker", http.StatusInternalServerError)
		return
	}

	utils.Tmpl.ExecuteTemplate(w, "workerprofile.html", worker)
}
