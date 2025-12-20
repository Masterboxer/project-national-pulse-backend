package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/models"
)

func GetTemplates(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, name, description, icon, created_at
			FROM templates
			ORDER BY id`)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var templates []models.Template
		for rows.Next() {
			var t models.Template
			if err := rows.Scan(
				&t.ID,
				&t.Name,
				&t.Description,
				&t.Icon,
				&t.CreatedAt,
			); err != nil {
				http.Error(w, "Error scanning templates", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			templates = append(templates, t)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating templates", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(templates)
	}
}

func GetTemplateByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var t models.Template
		err := db.QueryRow(`
			SELECT id, name, description, icon, created_at
			FROM templates
			WHERE id = $1`, id).
			Scan(
				&t.ID,
				&t.Name,
				&t.Description,
				&t.Icon,
				&t.CreatedAt,
			)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Template not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
}

func CreateTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t models.Template
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if t.Name == "" || t.Description == "" || t.Icon == "" {
			http.Error(w, "name, description, and icon are required", http.StatusBadRequest)
			return
		}

		err := db.QueryRow(`
			INSERT INTO templates (name, description, icon, created_at)
			VALUES ($1, $2, $3, NOW())
			RETURNING id, created_at`,
			t.Name,
			t.Description,
			t.Icon,
		).Scan(&t.ID, &t.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}

func UpdateTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var t models.Template
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if t.Name == "" || t.Description == "" || t.Icon == "" {
			http.Error(w, "name, description, and icon are required", http.StatusBadRequest)
			return
		}

		res, err := db.Exec(`
			UPDATE templates
			SET name = $1,
			    description = $2,
			    icon = $3
			WHERE id = $4`,
			t.Name,
			t.Description,
			t.Icon,
			id,
		)
		if err != nil {
			http.Error(w, "Database update failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to check update result", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		var updated models.Template
		err = db.QueryRow(`
			SELECT id, name, description, icon, created_at
			FROM templates
			WHERE id = $1`, id).
			Scan(
				&updated.ID,
				&updated.Name,
				&updated.Description,
				&updated.Icon,
				&updated.CreatedAt,
			)
		if err != nil {
			http.Error(w, "Failed to fetch updated template", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)
	}
}

func DeleteTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		res, err := db.Exec(`DELETE FROM templates WHERE id = $1`, id)
		if err != nil {
			http.Error(w, "Failed to delete template", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to check delete result", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Template deleted successfully",
		})
	}
}
