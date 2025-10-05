package routes

import (
	"net/http"

	"gocket/internal/api/handlers"
	"gocket/internal/api/middleware"
)

func SetupRoutes(articleHandler *handlers.ArticleHandler) http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			articleHandler.SaveArticle(w, r)
		case http.MethodGet:
			articleHandler.GetArticles(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/articles/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			articleHandler.GetArticle(w, r)
		case http.MethodPut:
			articleHandler.UpdateArticle(w, r)
		case http.MethodDelete:
			articleHandler.DeleteArticle(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Serve article reading page
	mux.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "./web/static/article.html")
	})

	// Serve static files (for web interface) - exclude /article path
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web/static/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	// Add middleware
	return middleware.CORS(middleware.Logging(mux))
}
