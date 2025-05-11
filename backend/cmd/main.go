package main

import (
	"backend/handler"
	"backend/pkg/store"
	"log"
	"net/http"
)

func addCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins or specify your frontend URL
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			return
		}

		h(w, r)
	}
}

func main() {
	err := store.InitDB("root:@tcp(127.0.0.1:3306)/user_authentication")
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	http.HandleFunc("/api/user/", addCORS(handler.SaveUserHandler))
	http.HandleFunc("/api/user/login", addCORS(handler.AuthenticateUserHandler))
	http.HandleFunc("/api/user/me", addCORS(handler.GetCurrentUserHandler))
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
