package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	heartbeatPath := "/ping"
	muxWithHeartbeat := HeartbeatMiddleware(mux, heartbeatPath)

	// Agrega ruta POST para /authenticate
	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		app.Authenticate(w, r)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(muxWithHeartbeat)

	return handler
}

func HeartbeatMiddleware(next http.Handler, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.Write([]byte("OK"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
