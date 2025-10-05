package main

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
	defer db.Close()

	// Ensure a tiny schema
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

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Trivial query
		var count int
		if err := db.QueryRow("SELECT COUNT(1) FROM articles").Scan(&count); err != nil {
			http.Error(w, fmt.Sprintf("db error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("Gocket up. Articles: %d\n", count)))
	})

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
