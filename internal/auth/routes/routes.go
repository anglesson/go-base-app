package routes

import (
	"github.com/anglesson/go-base-app/internal/auth/controllers"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router) {
	s := r.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/register", controllers.Register).Methods("POST")
	s.HandleFunc("/login", controllers.Login).Methods("POST")
	s.HandleFunc("/recover", controllers.RecoverPassword).Methods("POST")
	s.HandleFunc("/verify-reset-token", controllers.VerifyResetToken).Methods("POST")
	s.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")
}
