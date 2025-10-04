package main

import (
	"log"
	"net/http"

	"gocket/internal/api/handlers"
	"gocket/internal/api/routes"
	"gocket/internal/database"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize handlers
	articleHandler := handlers.NewArticleHandler(db)

	// Setup routes
	router := routes.SetupRoutes(articleHandler)

	// Start server
	log.Println("Starting Gocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
