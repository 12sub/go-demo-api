package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/12sub/websockets/auth"
	"github.com/12sub/websockets/db"
	"github.com/12sub/websockets/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// Hash password before starting
	if err := user.HashPassword(); err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User registered successfully"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	// Fetch user from DB
	result := db.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token, err := auth.GenerateJsonWebToken(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
