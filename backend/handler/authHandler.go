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

// SaveUserHandler handles the user registration
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

// AuthenticateUserHandler handles user login
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

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Refresh token generation failed", http.StatusInternalServerError)
		return
	}

	// Set the refresh token in the cookie (secure cookie)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	// Respond with the access token
	json.NewEncoder(w).Encode(map[string]string{
		"token": accessToken,
		"email": user.Username,
	})
}

// RefreshTokenHandler handles refreshing the access token using the refresh token
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	// Parse the refresh token
	claims, err := auth.ParseRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64) // user_id is stored as float64
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	user, err := store.GetUserByID(int(userID))
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	// Return the new access token
	json.NewEncoder(w).Encode(map[string]string{
		"token": accessToken,
	})
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Expire the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// GetCurrentUserHandler retrieves current user info using the access token
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

	// Respond with the current user info
	json.NewEncoder(w).Encode(map[string]string{
		"email": user.Username,
	})
}
