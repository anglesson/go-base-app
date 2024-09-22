package main

import (
	"log"
	"net/http"

	"github.com/anglesson/go-base-app/config"
	"github.com/anglesson/go-base-app/database"
	"github.com/anglesson/go-base-app/routes"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()
	database.Connect()

	r := mux.NewRouter()
	routes.RegisterAuthRoutes(r)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
