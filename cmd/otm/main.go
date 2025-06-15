package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"otm/internal/middleware"
	"otm/internal/routes"
	"otm/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	secret := os.Getenv("SECRET_KEY")
	port := os.Getenv("PORT")
	if len(secret) != 44 { // base64-encoded 32-byte key = 44 chars
		log.Fatal("SECRET_KEY must be a base64-encoded 32-byte key")
	}

	dbConn, err := storage.InitDB("otm.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	if err := storage.InitLogs(dbConn); err != nil {
		log.Fatal("Failed to initialize logs:", err)
	}

	db := &storage.DBHandle{Conn: dbConn}

	r := chi.NewRouter()
	r.Use(middleware.RateLimitMiddleware)
	routes.RegisterRoutes(r, db)

	log.Println(fmt.Sprintf("OTM server running at http://localhost:%s", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
