package handler

import (
	"backend/internal/auth"
	"backend/pkg/store"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type RegisterRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Hashing error", http.StatusInternalServerError)
		return
	}

	user := store.User{Username: req.Username, PasswordHash: hash}
	if err := store.SaveUser(user); err != nil {
		http.Error(w, "Could not save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	user, err := store.GetUserByUsername(creds.Username)
	if err != nil || user == nil || !auth.CheckPasswordHash(creds.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	// Optionally set refresh token as cookie (not implemented fully here)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "placeholder",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"email": user.Username,
	})
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := auth.ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	user, err := store.GetUserByUsername(username)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"email": user.Username,
	})
}
