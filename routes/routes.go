package routes

import (
	"github.com/anglesson/go-base-app/controllers"
	"github.com/anglesson/go-base-app/middlewares"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router) {
	r.Use(middlewares.CORS)
	r.Use(middlewares.JSON)
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/recover", controllers.RecoverPassword).Methods("POST")
	r.HandleFunc("/verify-reset-token", controllers.VerifyResetToken).Methods("POST")
	r.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")
}
