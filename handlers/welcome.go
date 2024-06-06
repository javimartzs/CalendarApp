package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"
	"net/http"
	"strconv"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {

	// Verificamos si el usuario esta logeado
	if !models.IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extraemos los valores del formulario de crear trabajadores
	if r.Method == http.MethodPost {
		workerId := r.FormValue("workerId")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		phone := r.FormValue("phone")
		store := r.FormValue("store")

		worker := models.Worker{
			Firstname: firstname,
			Lastname:  lastname,
			Phone:     phone,
			Store:     store,
		}

		if workerId != "" {
			id, err := strconv.Atoi(workerId)
			if err == nil {
				worker.ID = id
				if err := models.UpdateWorker(worker); err != nil {
					http.Error(w, "Unable to update worker", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid worker ID", http.StatusBadRequest)
				return
			}
		} else {
			if err := models.AddWorker(worker); err != nil {
				http.Error(w, "Unable to add worker", http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	// Proceso para visualizar la tabla de trabajadores
	workers, err := models.GetWorkers()
	if err != nil {
		http.Error(w, "Unable to retrieve workers", http.StatusInternalServerError)
		return
	}

	// Creamos una structura temporal para utilizar en este handler
	tmplData := struct {
		Username string
		Workers  []models.Worker
	}{
		Username: r.FormValue("username"),
		Workers:  workers,
	}

	utils.Tmpl.ExecuteTemplate(w, "welcome.html", tmplData)
}
