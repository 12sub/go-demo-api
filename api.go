package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"github.com/gorilla/websocket"
)

type APIServer struct {
	addr string
}

type Middleware func(http.Handler) http.HandlerFunc

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

var limiter = rate.NewLimiter(1, 5)
var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	for {
		_, msg, _ := conn.ReadMessage()
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}


func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: " + userID))
	})
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlechain := MiddlewareChain(RequestLoggerMiddleware, RequireAuthMiddleWare)
	server := http.Server{
		Addr:         s.addr,
		Handler:      middlechain(router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Starting Server %s ", s.addr)

	return server.ListenAndServe()
}

// Creating a Request Logger Middleware
func RequestLoggerMiddleware(rest http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path %s", r.Method, r.URL.Path)
		rest.ServeHTTP(w, r)
	}
}

// Middleware to authenticate users and routes
func RequireAuthMiddleWare(rest http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First step: check if users is authenticated
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// if authenticated
		rest.ServeHTTP(w, r)
	}
}

func MiddlewareChain(mw ...Middleware) Middleware {
	return func(h http.Handler) http.HandlerFunc {
		for i := len(mw) - 1; i >= 0; i-- {
			h = mw[i](h)
		}

		return h.ServeHTTP
	}
}
