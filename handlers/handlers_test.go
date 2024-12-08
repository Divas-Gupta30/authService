package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Divas-Gupta30/authService/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	handlers.InitInMemoryStore()
	router := mux.NewRouter()
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.HandleFunc("/anyprotectedroute", handlers.Protected).Methods("GET")
	router.HandleFunc("/refresh", handlers.Refresh).Methods("POST")
	router.HandleFunc("/revoke", handlers.RevokeToken).Methods("POST")
	return router
}

func Test_SignUp(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		user       map[string]string
		wantStatus int
	}{
		{"ValidUser", map[string]string{"email": "newuser@example.com", "password": "newpass"}, http.StatusOK},
		{"ExistingUser", map[string]string{"email": "test@example.com", "password": "testpass"}, http.StatusConflict},
		{"InvalidUser", map[string]string{"email": "invalid", "password": ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func Test_SignIn(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		user       map[string]string
		wantStatus int
	}{
		{"ValidUser", map[string]string{"email": "test@example.com", "password": "testpass"}, http.StatusOK},
		{"InvalidUser", map[string]string{"email": "test@example.com", "password": "wrongpass"}, http.StatusUnauthorized},
		{"NonExistentUser", map[string]string{"email": "nonexistent@example.com", "password": "pass"}, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func Test_Protected(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"InvalidToken", "invalid_token", http.StatusUnauthorized},
		{"ExpiredToken", "expired_token", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/anyprotectedroute", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func Test_Refresh(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"InvalidToken", "invalid_token", http.StatusUnauthorized},
		{"ExpiredToken", "expired_token", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/refresh", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func Test_RevokeToken(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{"InvalidToken", "invalid_token", http.StatusUnauthorized},
		{"NonExistentToken", "nonexistent_token", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/revoke", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
