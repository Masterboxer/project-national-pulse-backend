package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"masterboxer.com/project-mokuhyo/database"
	"masterboxer.com/project-mokuhyo/routes"
)

func main() {
	// Initialize database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Create router
	router := mux.NewRouter()

	// Create routes
	routes.CreateUserRoutes(db, router)

	// Wrap router with middleware
	handler := corsMiddleware(jsonContentTypeMiddleware(router))

	// Start server
	log.Println("Starting server on :8000...")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

// Middleware to set the content-type to JSON
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

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
