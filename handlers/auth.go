package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Divas-Gupta30/authService/helper"
	"github.com/Divas-Gupta30/authService/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
)

var (
	users       = make(map[string]models.User)
	tokens      = make(map[string]models.Token)
	tokenExpiry = time.Hour
)

func InitInMemoryStore() {
	users["test@example.com"] = models.User{Email: "test@example.com", Password: "testpass"}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate user using struct tags
	if err := validator.New().Struct(req); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := users[req.Email]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Add the user to the in-memory store
	users[req.Email] = req

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User signed up successfully"})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.User
	// Parse JSON payload into struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate user using struct tags
	if err := validator.New().Struct(req); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[req.Email]
	if !exists || user.Password != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := helper.GenerateJWT(req.Email, tokenExpiry)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Store the token in memory with expiration
	tokens[req.Email] = models.Token{
		UserID:  req.Email,
		Token:   tokenString,
		Expires: time.Now().Add(tokenExpiry).Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "message": "Signing Success"})
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	tokenString, err := helper.ExtractTokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Parse and validate the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &helper.Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*helper.Payload)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	storedToken, exists := tokens[claims.UserID]
	if !exists || storedToken.Expires < time.Now().Unix() {
		http.Error(w, "Token expired or not found", http.StatusUnauthorized)
		return
	}

	// Check if the token is about to expire (within 5 minutes)
	if storedToken.Expires-time.Now().Unix() < 5*60 {
		newTokenString, err := helper.GenerateJWT(claims.UserID, tokenExpiry)
		if err != nil {
			http.Error(w, "Could not generate new token", http.StatusInternalServerError)
			return
		}

		// Store the new token with updated expiration time
		tokens[claims.UserID] = models.Token{
			UserID:  claims.UserID,
			Token:   newTokenString,
			Expires: time.Now().Add(tokenExpiry).Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"new_token": newTokenString})
	} else {
		// Token is still valid
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Token is still valid"})
	}
}

func Protected(w http.ResponseWriter, r *http.Request) {
	tokenString, err := helper.ExtractTokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse and validate the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &helper.Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Validate if the token exists in the in-memory store
	claims, ok := token.Claims.(*helper.Payload)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	storedToken, exists := tokens[claims.UserID]
	if !exists || storedToken.Expires < time.Now().Unix() {
		http.Error(w, "Token expired or not found", http.StatusUnauthorized)
		return
	}

	// Access granted to the protected route
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Protected route accessed"})
}

func RevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := helper.ExtractTokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &helper.Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*helper.Payload)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if the token exists in the in-memory store
	if _, exists := tokens[claims.UserID]; exists {
		delete(tokens, claims.UserID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Token revoked successfully"})
		return
	}

	http.Error(w, "Token not found", http.StatusNotFound)
}
