package main

import (
	"net/http"

	"github.com/12sub/websockets/db"
	"github.com/12sub/websockets/handlers"
	"github.com/12sub/websockets/middleware"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	r := mux.NewRouter()

	// Public Routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected Route (Requires Token)
	r.HandleFunc("/profile", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to your profile!"))
	})).Methods("GET")

	server := NewAPIServer(":8080")
	server.Run()
}
