package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"masterboxer.com/project-micro-journal/models"
)

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, username, display_name, age, 
            gender, email, password, created_at FROM users`)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Age,
				&u.Gender, &u.Email, &u.Password, &u.CreatedAt); err != nil {
				http.Error(w, "Error scanning user data", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			u.Password = ""
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func GetUserById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.User
		err := db.QueryRow(`SELECT id, username, display_name, age, 
            gender, email, password, created_at FROM users WHERE id = $1`, id).
			Scan(&u.ID, &u.Username, &u.DisplayName, &u.Age, &u.Gender, &u.Email,
				&u.Password, &u.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		u.Password = ""
		json.NewEncoder(w).Encode(u)
	}
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u models.User
		err := db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).
			Scan(&u.ID, &u.Username, &u.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database query failed", http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if u.Username == "" || u.DisplayName == "" || u.Email == "" || u.Password == "" {
			http.Error(w, "Username, display_name, email, and password are required", http.StatusBadRequest)
			return
		}

		if u.Age == 0 {
			http.Error(w, "Age is required and must be greater than 0", http.StatusBadRequest)
			return
		}

		if u.Gender == "" {
			http.Error(w, "Gender is required", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		err = db.QueryRow(
			`INSERT INTO users (username, display_name, age, gender, email, password, created_at) 
            VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id, created_at`,
			u.Username, u.DisplayName, u.Age, u.Gender, u.Email, string(hashedPassword),
		).Scan(&u.ID, &u.CreatedAt)

		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		u.Password = ""
		json.NewEncoder(w).Encode(u)
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		setClauses := []string{}
		args := []interface{}{}
		i := 1

		if u.Username != "" {
			setClauses = append(setClauses, "username = $"+strconv.Itoa(i))
			args = append(args, u.Username)
			i++
		}
		if u.DisplayName != "" {
			setClauses = append(setClauses, "display_name = $"+strconv.Itoa(i))
			args = append(args, u.DisplayName)
			i++
		}
		if u.Email != "" {
			setClauses = append(setClauses, "email = $"+strconv.Itoa(i))
			args = append(args, u.Email)
			i++
		}
		if u.Age != 0 {
			setClauses = append(setClauses, "age = $"+strconv.Itoa(i))
			args = append(args, u.Age)
			i++
		}
		if u.Gender != "" {
			setClauses = append(setClauses, "gender = $"+strconv.Itoa(i))
			args = append(args, u.Gender)
			i++
		}

		if len(setClauses) == 0 {
			http.Error(w, "No fields provided for update", http.StatusBadRequest)
			return
		}

		sqlStr := "UPDATE users SET " + strings.Join(setClauses, ", ") +
			" WHERE id = $" + strconv.Itoa(i)
		args = append(args, id)

		_, err := db.Exec(sqlStr, args...)
		if err != nil {
			http.Error(w, "Database update failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var updatedUser models.User
		err = db.QueryRow(`SELECT id, username, display_name, age, 
            gender, email, password, created_at FROM users WHERE id = $1`, id).
			Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.DisplayName,
				&updatedUser.Age, &updatedUser.Gender, &updatedUser.Email,
				&updatedUser.Password, &updatedUser.CreatedAt)

		if err != nil {
			http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		updatedUser.Password = ""
		json.NewEncoder(w).Encode(updatedUser)
	}
}
