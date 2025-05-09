package main

import (
	"backend/handler"
	"backend/pkg/store"
	"log"
	"net/http"
)

func main() {
	err := store.InitDB("root:@tcp(127.0.0.1:3306)/user_authentication")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	http.HandleFunc("/api/user/", handler.SaveUserHandler)
	// http.HandleFunc("/api/user/", handler.GetUserHandler)
	http.HandleFunc("/api/user/login", handler.AuthenticateUserHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
