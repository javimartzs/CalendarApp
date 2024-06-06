package handlers

import (
	"CalendarApp/models"
	"CalendarApp/utils"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Aquí deberías validar el usuario y la contraseña
	if models.IsValidUser(username, password) {
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: username,
			Path:  "/",
		})
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	} else {
		utils.Tmpl.ExecuteTemplate(w, "login.html", "Usuario o contraseña incorrectos")
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
