package main

import (
	"fmt"
	"github.com/Divas-Gupta30/authService/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize in-memory user and token maps (in place of a database)
	handlers.InitInMemoryStore()

	// HTTP Routes
	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/signin", handlers.SignIn)
	http.HandleFunc("/protected", handlers.Protected)
	http.HandleFunc("/refresh", handlers.Refresh)

	// Start the server
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
