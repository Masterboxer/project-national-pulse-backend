package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/handlers"
)

func CreateNotificationRoutes(db *sql.DB, router *mux.Router) *mux.Router {

	router.HandleFunc("/fcm/register-token", handlers.RegisterFCMToken(db)).Methods("POST")

	return router
}
