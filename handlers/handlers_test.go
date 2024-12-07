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
	router.HandleFunc("/protected", handlers.Protected).Methods("GET")
	router.HandleFunc("/refresh", handlers.Refresh).Methods("POST")
	router.HandleFunc("/revoke", handlers.RevokeToken).Methods("POST")
	return router
}

func Test_SignUp(t *testing.T) {
	router := setupRouter()

	user := map[string]string{"email": "newuser@example.com", "password": "newpass"}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_SignIn(t *testing.T) {
	router := setupRouter()

	user := map[string]string{"email": "test@example.com", "password": "testpass"}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_Protected(t *testing.T) {
	router := setupRouter()

	// First, sign in to get a token
	user := map[string]string{"email": "test@example.com", "password": "testpass"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	token := response["token"]

	// Use the token to access the protected route
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_Refresh(t *testing.T) {
	router := setupRouter()

	// First, sign in to get a token
	user := map[string]string{"email": "test@example.com", "password": "testpass"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	token := response["token"]

	// Use the token to refresh
	req, _ = http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRevokeToken(t *testing.T) {
	router := setupRouter()

	// First, sign in to get a token
	user := map[string]string{"email": "test@example.com", "password": "testpass"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	token := response["token"]

	// Use the token to revoke
	req, _ = http.NewRequest("POST", "/revoke", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
