package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	// Middleware para CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Ruta para Broker
	brokerHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			app.Broker(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// Ruta para HandleSubmission
	handleSubmissionHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			app.HandleSubmission(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	mux.HandleFunc("/", brokerHandler)
	mux.HandleFunc("/handle", handleSubmissionHandler)

	// Aplicando los middlewares
	handlerWithCors := c.Handler(mux)
	handlerWithHeartbeat := HeartbeatMiddleware(handlerWithCors.ServeHTTP)

	return http.HandlerFunc(handlerWithHeartbeat)
}

func HeartbeatMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			// Aquí puedes colocar la lógica de tu middleware Heartbeat.
			// Por ahora, sólo lo simulo con un simple Write.
			w.Write([]byte("OK"))
			return
		}
		next(w, r)
	}
}
