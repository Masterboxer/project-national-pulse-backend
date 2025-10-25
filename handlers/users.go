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
	"masterboxer.com/project-mokuhyo/models"
)

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, COALESCE(age, 0), COALESCE(gender, ''), email, password FROM users")
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Gender, &u.Email, &u.Password); err != nil {
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
		err := db.QueryRow("SELECT id, name, COALESCE(age, 0), COALESCE(gender, ''), email, password FROM users WHERE id = $1", id).
			Scan(&u.ID, &u.Name, &u.Age, &u.Gender, &u.Email, &u.Password)
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
		err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).
			Scan(&u.ID, &u.Name, &u.Email)
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
		json.NewDecoder(r.Body).Decode(&u)

		if u.Password == "" {
			http.Error(w, "Password cannot be empty", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		err = db.QueryRow(
			"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
			u.Name, u.Email, string(hashedPassword),
		).Scan(&u.ID)

		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
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

		if u.Name != "" {
			setClauses = append(setClauses, "name = $"+strconv.Itoa(i))
			args = append(args, u.Name)
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

		sqlStr := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(i)
		args = append(args, id)

		_, err := db.Exec(sqlStr, args...)
		if err != nil {
			http.Error(w, "Database update failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var updatedUser models.User
		err = db.QueryRow("SELECT id, name, COALESCE(age, 0), COALESCE(gender, ''), email, password FROM users WHERE id = $1", id).
			Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Age, &updatedUser.Gender, &updatedUser.Email, &updatedUser.Password)

		if err != nil {
			http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		updatedUser.Password = ""
		json.NewEncoder(w).Encode(updatedUser)
	}
}
