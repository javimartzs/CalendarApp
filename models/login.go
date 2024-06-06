package models

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type User struct {
	Username string
	Password string
}
type Config struct {
	Users []User
}

var config Config

// Funcion para importar y leer el fichero config
func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error import Config file; %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error parsing Config file: %v", err)
	}
}

// Funcion para validar usuario y contraseña
func IsValidUser(username, password string) bool {
	for _, user := range config.Users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

// Funcion para comprobar si el usuario esta logeado
func IsLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		return false
	}
	return true
}
