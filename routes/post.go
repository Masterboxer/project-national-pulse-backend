package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/handlers"
)

func CreatePostRoutes(db *sql.DB, router *mux.Router) *mux.Router {
	router.HandleFunc("/posts", handlers.GetPosts(db)).Methods("GET")
	router.HandleFunc("/posts/today", handlers.GetTodayPostForUser(db)).Methods("GET")
	router.HandleFunc("/posts", handlers.CreatePost(db)).Methods("POST")

	router.HandleFunc("/posts/{id}", handlers.GetPostByID(db)).Methods("GET")
	router.HandleFunc("/posts/{id}", handlers.DeletePost(db)).Methods("DELETE")

	return router
}
