package models

import (
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type Worker struct {
	ID           int
	Firstname    string
	Lastname     string
	Phone        string
	Store        string
	DisplayOrder int
}

func AddWorker(worker Worker) error {
	// Inserta el nuevo trabajador con un display_order temporal
	insertQuery := `INSERT INTO workers (firstname, lastname, phone, store, display_order) VALUES (?, ?, ?, ?, 9999)`
	_, err := DB.Exec(insertQuery, worker.Firstname, worker.Lastname, worker.Phone, worker.Store)
	if err != nil {
		return err
	}

	// Reordena todos los trabajadores
	return ReorderWorkers()
}

func UpdateWorker(worker Worker) error {
	// Actualiza el trabajador existente
	updateQuery := `UPDATE workers SET firstname = ?, lastname = ?, phone = ?, store = ? WHERE id = ?`
	_, err := DB.Exec(updateQuery, worker.Firstname, worker.Lastname, worker.Phone, worker.Store, worker.ID)
	if err != nil {
		return err
	}

	// Reordena todos los trabajadores
	return ReorderWorkers()
}

func ReorderWorkers() error {
	// Obtiene todos los trabajadores
	workers, err := GetWorkers()
	if err != nil {
		return err
	}

	// Ordena los trabajadores por tienda y luego por nombre
	sort.Slice(workers, func(i, j int) bool {
		if workers[i].Store == workers[j].Store {
			return workers[i].Firstname < workers[j].Firstname
		}
		return workers[i].Store < workers[j].Store
	})

	// Actualiza el campo display_order de cada trabajador
	for i, worker := range workers {
		updateOrderQuery := `UPDATE workers SET display_order = ? WHERE id = ?`
		_, err := DB.Exec(updateOrderQuery, i+1, worker.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetWorkers() ([]Worker, error) {
	query := `SELECT id, firstname, lastname, phone, store, display_order FROM workers ORDER BY display_order`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workers := []Worker{}
	for rows.Next() {
		var worker Worker
		if err := rows.Scan(&worker.ID, &worker.Firstname, &worker.Lastname, &worker.Phone, &worker.Store, &worker.DisplayOrder); err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}
	return workers, nil
}

func GetWorkerByID(id int) (Worker, error) {
	query := `SELECT id, firstname, lastname, phone, store, display_order FROM workers WHERE id = ?`
	row := DB.QueryRow(query, id)

	var worker Worker
	if err := row.Scan(&worker.ID, &worker.Firstname, &worker.Lastname, &worker.Phone, &worker.Store, &worker.DisplayOrder); err != nil {
		return worker, err
	}
	return worker, nil
}

func DeleteWorker(id int) error {
	query := `DELETE FROM workers WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Reordena los trabajadores después de eliminar uno
	return ReorderWorkers()
}
