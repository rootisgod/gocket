package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and creates tables
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./gocket.db")
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	articlesTable := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		excerpt TEXT,
		author TEXT,
		published_date DATETIME,
		saved_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		read_status BOOLEAN DEFAULT FALSE,
		tags TEXT, -- JSON array of tags
		word_count INTEGER,
		reading_time INTEGER, -- estimated minutes
		thumbnail_url TEXT,
		domain TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_articles_url ON articles(url);",
		"CREATE INDEX IF NOT EXISTS idx_articles_saved_date ON articles(saved_date);",
		"CREATE INDEX IF NOT EXISTS idx_articles_read_status ON articles(read_status);",
		"CREATE INDEX IF NOT EXISTS idx_articles_domain ON articles(domain);",
	}

	// Execute table creation
	if _, err := db.Exec(articlesTable); err != nil {
		return err
	}

	// Execute index creation
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return err
		}
	}

	return nil
}
