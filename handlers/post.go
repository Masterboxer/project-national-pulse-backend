package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/models"
	"masterboxer.com/project-micro-journal/services"
)

func GetPostsByUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr, ok := vars["userId"]
		if !ok || userIDStr == "" {
			http.Error(w, "userId parameter missing", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT id, user_id, template_id, text, 
			       COALESCE(photo_path, '') as photo_path, 
			       created_at
			FROM posts
			WHERE user_id = $1
			ORDER BY created_at DESC`,
			userID)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Printf("GetPostsByUser error: %v", err)
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
				log.Printf("GetPostsByUser scan error: %v", err)
				return
			}
			posts = append(posts, p)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating posts", http.StatusInternalServerError)
			log.Printf("GetPostsByUser rows error: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
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
			http.Error(w, "template_id is required", http.StatusBadRequest)
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

		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		var postCount int
		err := db.QueryRow(`
			SELECT COUNT(*) 
			FROM posts 
			WHERE user_id = $1 
			  AND created_at >= $2 
			  AND created_at < $3`,
			p.UserID, startOfDay, endOfDay).Scan(&postCount)
		if err != nil {
			http.Error(w, "Failed to check daily limit", http.StatusInternalServerError)
			log.Println("CreatePost daily limit check error:", err)
			return
		}

		if postCount > 0 {
			http.Error(w, "Daily post limit reached (1 post per day)", http.StatusForbidden)
			return
		}

		err = db.QueryRow(`
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

		go notifyBuddiesOfNewPost(db, p.UserID, p.Text)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

func notifyBuddiesOfNewPost(db *sql.DB, userID int, postText string) {
	var displayName string
	err := db.QueryRow(`SELECT display_name FROM users WHERE id = $1`, userID).Scan(&displayName)
	if err != nil {
		log.Printf("Error fetching user display name for notifications: %v", err)
		displayName = "A friend"
	}

	rows, err := db.Query(`
		SELECT DISTINCT ft.token
		FROM buddies b
		JOIN fcm_tokens ft ON b.buddy_id = ft.user_id
		WHERE b.user_id = $1 
		  AND ft.token IS NOT NULL 
		  AND ft.token != ''`,
		userID)
	if err != nil {
		log.Printf("Error fetching buddy FCM tokens: %v", err)
		return
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scanning FCM token: %v", err)
			continue
		}
		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		log.Printf("No FCM tokens found for user %d's buddies", userID)
		return
	}

	title := fmt.Sprintf("%s posted today!", displayName)
	body := postText
	if len(body) > 100 {
		body = body[:97] + "..."
	}

	data := map[string]string{
		"type":    "new_post",
		"user_id": strconv.Itoa(userID),
	}

	successCount, failureCount, err := services.SendMultipleNotifications(tokens, title, body, data)
	if err != nil {
		log.Printf("Error sending notifications to buddies: %v", err)
		return
	}

	log.Printf("Sent notifications for new post by user %d: %d successful, %d failed",
		userID, successCount, failureCount)
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

func GetBuddyPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr, ok := vars["userId"]
		if !ok || userIDStr == "" {
			http.Error(w, "userId parameter missing", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}

		thirtysixHoursAgo := time.Now().Add(-36 * time.Hour)

		rows, err := db.Query(`
            SELECT p.id, p.user_id, p.template_id, p.text, 
                   COALESCE(p.photo_path, '') as photo_path, 
                   p.created_at, 
                   u.username, u.display_name
            FROM posts p
            JOIN users u ON p.user_id = u.id
            WHERE p.user_id = $1
              AND p.created_at >= $2
            ORDER BY p.created_at DESC`,
			userID, thirtysixHoursAgo)
		if err != nil {
			http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
			log.Println("GetBuddyPosts user posts error:", err)
			return
		}
		defer rows.Close()

		var userPosts []models.PostWithUser
		for rows.Next() {
			var p models.PostWithUser
			if err := rows.Scan(
				&p.ID,
				&p.UserID,
				&p.TemplateID,
				&p.Text,
				&p.PhotoPath,
				&p.CreatedAt,
				&p.Username,
				&p.DisplayName,
			); err != nil {
				http.Error(w, "Error scanning user posts", http.StatusInternalServerError)
				log.Println("GetBuddyPosts user scan error:", err)
				return
			}
			userPosts = append(userPosts, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating user posts", http.StatusInternalServerError)
			log.Println("GetBuddyPosts user rows error:", err)
			return
		}

		rows, err = db.Query(`
            SELECT
                p.id,
                p.user_id,
                p.template_id,
                p.text,
                COALESCE(p.photo_path, '') as photo_path,
                p.created_at,
                u.username,
                u.display_name
            FROM posts p
            JOIN buddies b ON p.user_id = b.buddy_id
            JOIN users u ON p.user_id = u.id
            WHERE b.user_id = $1
              AND p.user_id != $2
              AND p.created_at >= $3
            ORDER BY p.created_at DESC
            LIMIT 49`,
			userID, userID, thirtysixHoursAgo)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println("GetBuddyPosts buddy posts error:", err)
			return
		}
		defer rows.Close()

		var buddyPosts []models.PostWithUser
		for rows.Next() {
			var p models.PostWithUser
			if err := rows.Scan(
				&p.ID,
				&p.UserID,
				&p.TemplateID,
				&p.Text,
				&p.PhotoPath,
				&p.CreatedAt,
				&p.Username,
				&p.DisplayName,
			); err != nil {
				http.Error(w, "Error scanning buddy posts", http.StatusInternalServerError)
				log.Println("GetBuddyPosts buddy scan error:", err)
				return
			}
			buddyPosts = append(buddyPosts, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating buddy posts", http.StatusInternalServerError)
			log.Println("GetBuddyPosts buddy rows error:", err)
			return
		}

		var feed []models.PostWithUser
		feed = append(feed, userPosts...)
		feed = append(feed, buddyPosts...)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feed)
	}
}
