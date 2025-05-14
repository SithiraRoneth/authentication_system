package main

import (
	"backend/handler"
	"backend/pkg/store"
	"fmt"
	"log"
	"net/http"
)

func addCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

func main() {
	dbUser := "root"
	dbPass := "1234"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "user_authentication"

	// Create the Data Source Name (DSN) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize the database connection
	err := store.InitDB(dsn)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	// Set up the routes with CORS handling
	http.HandleFunc("/api/user/", addCORS(handler.SaveUserHandler))
	http.HandleFunc("/api/user/login", addCORS(handler.AuthenticateUserHandler))
	http.HandleFunc("/api/user/me", addCORS(handler.GetCurrentUserHandler))
	http.HandleFunc("/api/user/refresh", addCORS(handler.RefreshTokenHandler))
	http.HandleFunc("/api/user/logout", addCORS(handler.LogoutHandler))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
