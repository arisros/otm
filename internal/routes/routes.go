package routes

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"otm/internal/storage"
)

// RegisterRoutes injects DB and mounts routes
func RegisterRoutes(r chi.Router, db *storage.DBHandle) {
	// Serve homepage form
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "index.html"))
	})

	// Serve read message HTML
	r.Get("/msg/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") == "application/json" {
			handleReadMessage(db).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, filepath.Join("web", "templates", "message.html"))
		}
	})

	// API endpoints
	r.Post("/api/messages", handleCreateMessage(db))
	r.Get("/api/msg/{id}", handleReadMessage(db)) // Optional direct API read
}
