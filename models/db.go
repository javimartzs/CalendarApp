package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createWorkersTableQuery := `
    CREATE TABLE IF NOT EXISTS workers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        firstname TEXT,
        lastname TEXT,
        phone TEXT,
        store TEXT,
        display_order INTEGER
    );`
	if _, err := DB.Exec(createWorkersTableQuery); err != nil {
		log.Fatalf("Error creating workers table: %v", err)
	}
}
