package main

import (
	"CalendarApp/handlers"
	"CalendarApp/models"
	"CalendarApp/utils"
	"log"
	"net/http"
)

func main() {
	models.LoadConfig()
	models.InitDB("workers.db")

	// Ensure templates are initialized with the required functions
	utils.InitTemplates()

	http.HandleFunc("/", handlers.LoginHandler)
	http.HandleFunc("/welcome", handlers.WelcomeHandler)
	http.HandleFunc("/worker/", handlers.WorkerProfileHandler)
	http.HandleFunc("/calendar", handlers.CalendarHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
