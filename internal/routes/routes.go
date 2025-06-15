package routes

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"otm/internal/storage"
)

var templates = map[string]string{
	"home":    "index.html",
	"message": "message.html",
}

func RegisterRoutes(r chi.Router, db *storage.DBHandle) {
	// static html
	r.Get("/", serveTemplate("home"))
	r.Get("/msg/{id}", serveTemplate("message"))

	// rest api
	r.Post("/api/messages", handleCreateMessage(db))
	r.Get("/api/msg/{id}", handleReadMessage(db))

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

}

func serveTemplate(name string) http.HandlerFunc {
	filename, ok := templates[name]
	if !ok {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Template not found", http.StatusNotFound)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", filename))
	}
}
