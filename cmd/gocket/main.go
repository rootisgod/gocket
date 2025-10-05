package main

// Gocket: minimal HTTP service backed by SQLite to track saved articles.

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	// Open SQLite DB with default settings
	db, err := sql.Open("sqlite", "file:gocket.db")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	// Always close the database connection on shutdown
	defer db.Close()

	// Ensure a tiny schema for storing saved articles
	if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS articles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT NOT NULL,
            title TEXT,
            saved_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Configure HTTP routes and handlers
	mux := http.NewServeMux()
	// Health check endpoint for liveness/readiness probes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
	})
	// Root endpoint: returns a simple status and article count
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Query the total number of saved articles
		var count int
		if err := db.QueryRow("SELECT COUNT(1) FROM articles").Scan(&count); err != nil {
			http.Error(w, fmt.Sprintf("db error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("Gocket up. Articles: %d\n", count)))
	})

	// Start the HTTP server
	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
