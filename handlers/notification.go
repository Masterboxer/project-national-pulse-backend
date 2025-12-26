package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type TokenRequest struct {
	Token     string `json:"token"`
	UserID    int    `json:"user_id"`
	Timestamp string `json:"timestamp,omitempty"`
}

func RegisterTokenHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req TokenRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("JSON decode error: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 {
			http.Error(w, "Valid user_id is required", http.StatusBadRequest)
			return
		}

		query := `UPDATE users SET fcm_token = $1, updated_at = NOW() WHERE id = $2`
		result, err := db.Exec(query, req.Token, req.UserID)

		if err != nil {
			log.Printf("Database error updating FCM token: %v", err)
			http.Error(w, "Failed to register token", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error checking rows affected: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Token registered successfully",
		})
	}
}
