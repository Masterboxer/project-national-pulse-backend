package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"masterboxer.com/project-micro-journal/database"
	"masterboxer.com/project-micro-journal/routes"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	routes.CreateUserRoutes(db, router)
	routes.CreateAuthenticationRoutes(db, router)
	routes.CreatePostRoutes(db, router)
	routes.CreateTemplateRoutes(db, router)

	handler := corsMiddleware(jsonContentTypeMiddleware(router))

	log.Println("Starting server on :8200...")
	log.Fatal(http.ListenAndServe(":8200", handler))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
