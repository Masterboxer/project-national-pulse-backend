package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/handlers"
)

func CreateTemplateRoutes(db *sql.DB, router *mux.Router) *mux.Router {

	router.HandleFunc("/templates", handlers.GetTemplates(db)).Methods("GET")
	router.HandleFunc("/templates/{id}", handlers.GetTemplateByID(db)).Methods("GET")
	router.HandleFunc("/templates", handlers.CreateTemplate(db)).Methods("POST")
	router.HandleFunc("/templates/{id}", handlers.UpdateTemplate(db)).Methods("PUT")
	router.HandleFunc("/templates/{id}", handlers.DeleteTemplate(db)).Methods("DELETE")

	return router
}
