package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/models"
)

func GetPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, user_id, template_id, text, photo_path, created_at
			FROM posts
			ORDER BY created_at DESC`)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var p models.Post
			if err := rows.Scan(
				&p.ID,
				&p.UserID,
				&p.TemplateID,
				&p.Text,
				&p.PhotoPath,
				&p.CreatedAt,
			); err != nil {
				http.Error(w, "Error scanning posts", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			posts = append(posts, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating posts", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func GetPostByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok || idStr == "" {
			http.Error(w, "ID parameter missing", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var p models.Post
		err = db.QueryRow(`
			SELECT id, user_id, template_id, text, photo_path, created_at
			FROM posts WHERE id = $1`, id).
			Scan(&p.ID, &p.UserID, &p.TemplateID, &p.Text, &p.PhotoPath, &p.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				log.Printf("GetPostByID error for id=%s: %v", idStr, err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p models.Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if p.UserID == 0 {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		if p.TemplateID == 0 {
			http.Error(w, "templateId is required", http.StatusBadRequest)
			return
		}
		if p.Text == "" {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}
		if len(p.Text) > 280 {
			http.Error(w, "text must be at most 280 characters", http.StatusBadRequest)
			return
		}

		err := db.QueryRow(`
			INSERT INTO posts (user_id, template_id, text, photo_path, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			RETURNING id, user_id, template_id, text, photo_path, created_at`,
			p.UserID,
			p.TemplateID,
			p.Text,
			p.PhotoPath,
		).Scan(
			&p.ID,
			&p.UserID,
			&p.TemplateID,
			&p.Text,
			&p.PhotoPath,
			&p.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			log.Println("CreatePost error:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var exists bool
		err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1)`, id).
			Scan(&exists)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if !exists {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(`DELETE FROM posts WHERE id = $1`, id)
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Post deleted successfully",
		})
	}
}

func GetTodayPostForUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int
		var err error

		vars := mux.Vars(r)
		if uidStr, ok := vars["userId"]; ok {
			userID, err = strconv.Atoi(uidStr)
			if err != nil {
				http.Error(w, "Invalid userId", http.StatusBadRequest)
				return
			}
		} else {
			uidStr := r.URL.Query().Get("user_id")
			if uidStr == "" {
				http.Error(w, "user_id is required", http.StatusBadRequest)
				return
			}
			userID, err = strconv.Atoi(uidStr)
			if err != nil {
				http.Error(w, "Invalid user_id", http.StatusBadRequest)
				return
			}
		}

		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var p models.Post
		err = db.QueryRow(`
			SELECT id, user_id, template_id, text, photo_path, created_at
			FROM posts
			WHERE user_id = $1
			  AND created_at >= $2
			  AND created_at < $3
			ORDER BY created_at DESC
			LIMIT 1`,
			userID, startOfDay, endOfDay,
		).Scan(
			&p.ID,
			&p.UserID,
			&p.TemplateID,
			&p.Text,
			&p.PhotoPath,
			&p.CreatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}
