package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anglesson/go-base-app/database"
	"github.com/anglesson/go-base-app/models"
	"github.com/anglesson/go-base-app/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	database.DB.First(&user, "email = ?", user.Email)

	if user.ID != 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	database.DB.Create(&user)

	token, _ := utils.GenerateJWT(user.Email)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	database.DB.Where("email = ?", input.Email).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateJWT(user.Email)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Gera o token de recuperação e define a expiração
	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(1 * time.Hour)

	// Atualiza o usuário com o token e a expiração
	user.ResetToken = token
	user.TokenExpiration = &expiration
	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	// Enviar e-mail de recuperação
	err = utils.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password reset email sent"))
}

func VerifyResetToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	if time.Now().After(*user.TokenExpiration) {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token is valid"))
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Password == "" || req.Token == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		http.Error(w, "Invalid token", http.StatusNotFound)
		return
	}

	if time.Now().After(*user.TokenExpiration) {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return
	}

	// Criptografa a nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.TokenExpiration = nil

	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password has been reset"))
}

// Generate random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
