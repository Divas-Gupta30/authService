package main

import (
	"fmt"
	"github.com/Divas-Gupta30/authService/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Initialize in-memory user and token maps (in place of a database)
	handlers.InitInMemoryStore()
	router := mux.NewRouter()

	// HTTP Routes
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.HandleFunc("/anyprotectedroute", handlers.Protected).Methods("GET")
	router.HandleFunc("/refresh", handlers.Refresh).Methods("POST")
	// for internal use
	router.HandleFunc("/revoke", handlers.RevokeToken).Methods("POST")

	// Start the server
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
