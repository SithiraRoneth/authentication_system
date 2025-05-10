package main

import (
	"backend/handler"
	"backend/pkg/store"
	"log"
	"net/http"
)

func addCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	err := store.InitDB("root:@tcp(127.0.0.1:3306)/user_authentication")
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	http.HandleFunc("/api/user/", addCORS(handler.SaveUserHandler))
	http.HandleFunc("/api/user/login", addCORS(handler.AuthenticateUserHandler))
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
