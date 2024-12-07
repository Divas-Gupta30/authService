package handlers

import (
	"net/http"
	"time"

	"github.com/Divas-Gupta30/authService/helper"
	"github.com/Divas-Gupta30/authService/models"
	"github.com/dgrijalva/jwt-go"
)

var users = make(map[string]models.User)
var tokens = make(map[string]models.Token)

func InitInMemoryStore() {
	users["test@example.com"] = models.User{Email: "test@example.com", Password: "testpass"}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if _, exists := users[email]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user := models.User{Email: email, Password: password}
	users[email] = user

	w.Write([]byte("User signed up successfully"))
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, exists := users[email]
	if !exists || user.Password != password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := helper.GenerateJWT(email)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	tokens[email] = models.Token{
		UserID:  email,
		Token:   tokenString,
		Expires: time.Now().Add(time.Hour).Unix(),
	}

	w.Write([]byte("Token: " + tokenString))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &helper.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*helper.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	storedToken, exists := tokens[claims.UserID]
	if !exists || storedToken.Expires < time.Now().Unix() {
		http.Error(w, "Token expired or not found", http.StatusUnauthorized)
		return
	}

	// Check if token is about to expire (within 5 minutes)
	if storedToken.Expires-time.Now().Unix() < 5*60 {
		newTokenString, err := helper.GenerateJWT(claims.UserID)
		if err != nil {
			http.Error(w, "Could not generate new token", http.StatusInternalServerError)
			return
		}

		tokens[claims.UserID] = models.Token{
			UserID:  claims.UserID,
			Token:   newTokenString,
			Expires: time.Now().Add(time.Hour).Unix(),
		}
		w.Write([]byte("New Token: " + newTokenString))
	} else {
		w.Write([]byte("Token is still valid"))
	}
}

func Protected(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, &helper.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Validate if the token exists in the in-memory store
	claims, ok := token.Claims.(*helper.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	storedToken, exists := tokens[claims.UserID]
	if !exists || storedToken.Expires < time.Now().Unix() {
		http.Error(w, "Token expired or not found", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Protected route accessed"))
}
