package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"Uebungsprojekt/backend/internal"
)

// --- main: Connect server to database and start server ---
func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL env var not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := internal.NewDB(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	srv, err := internal.NewServer(db, "templates/*.html")
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	log.Println("🚀 listening on :80")
	log.Fatal(http.ListenAndServe(":80", srv.Routes()))
}
