package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"masterboxer.com/project-national-pulse/handlers"
)

func CreateUserRoutes(db *sql.DB, router *mux.Router) *mux.Router {

	router.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.GetUserById(db)).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	return router
}
