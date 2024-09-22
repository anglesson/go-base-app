package main

import (
	"log"
	"net/http"

	"github.com/anglesson/go-base-app/internal/auth/routes"
	"github.com/anglesson/go-base-app/internal/config"
	"github.com/anglesson/go-base-app/internal/database"
	"github.com/anglesson/go-base-app/internal/middlewares"
	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	database.Connect()

	// Create a new router
	r := mux.NewRouter()

	// Apply the CORS and JSON middlewares to all routes
	r.Use(middlewares.CORS)
	r.Use(middlewares.JSON)

	// Register routes
	routes.RegisterAuthRoutes(r)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
