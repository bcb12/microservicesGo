package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *Config) routes() http.Handler {
	mux := http.HandlerFunc(app.Broker)

	heartbeatPath := "/ping"
	muxWithHeartbeat := HeartbeatMiddleware(mux, heartbeatPath)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(muxWithHeartbeat)

	http.Handle("/", handler)

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
