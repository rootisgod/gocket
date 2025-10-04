package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gocket/internal/database/repository"
	"gocket/internal/models"
	"gocket/internal/scraper"
)

type ArticleHandler struct {
	repo             *repository.ArticleRepository
	contentProcessor *scraper.ContentProcessor
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
	return &ArticleHandler{
		repo:             repository.NewArticleRepository(db),
		contentProcessor: scraper.NewContentProcessor(),
	}
}

// SaveArticle handles POST /api/articles
func (h *ArticleHandler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	var req models.SaveArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process the URL and extract content
	article, err := h.contentProcessor.ProcessURL(req.URL)
	if err != nil {
		http.Error(w, "Failed to process URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamps
	now := time.Now()
	article.SavedDate = now
	article.CreatedAt = now
	article.UpdatedAt = now

	// Save to database
	if err := h.repo.CreateArticle(article); err != nil {
		http.Error(w, "Failed to save article: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// GetArticles handles GET /api/articles
func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	perPage := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	articles, total, err := h.repo.GetAllArticles(page, perPage)
	if err != nil {
		http.Error(w, "Failed to retrieve articles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		PerPage:  perPage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetArticle handles GET /api/articles/:id
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Path[len("/api/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.repo.GetArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Article not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve article: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// UpdateArticle handles PUT /api/articles/:id
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Path[len("/api/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	article.ID = id
	article.UpdatedAt = time.Now()

	if err := h.repo.UpdateArticle(&article); err != nil {
		http.Error(w, "Failed to update article: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// DeleteArticle handles DELETE /api/articles/:id
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.URL.Path[len("/api/articles/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteArticle(id); err != nil {
		http.Error(w, "Failed to delete article: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
