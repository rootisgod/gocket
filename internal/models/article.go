package models

import "time"

// Article represents a saved web article
type Article struct {
	ID            int        `json:"id" db:"id"`
	URL           string     `json:"url" db:"url"`
	Title         string     `json:"title" db:"title"`
	Content       string     `json:"content" db:"content"`
	Excerpt       string     `json:"excerpt" db:"excerpt"`
	Author        string     `json:"author" db:"author"`
	PublishedDate *time.Time `json:"published_date" db:"published_date"`
	SavedDate     time.Time  `json:"saved_date" db:"saved_date"`
	ReadStatus    bool       `json:"read_status" db:"read_status"`
	Tags          []string   `json:"tags" db:"tags"`
	WordCount     int        `json:"word_count" db:"word_count"`
	ReadingTime   int        `json:"reading_time" db:"reading_time"`
	ThumbnailURL  string     `json:"thumbnail_url" db:"thumbnail_url"`
	Domain        string     `json:"domain" db:"domain"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// SaveArticleRequest represents the request to save a new article
type SaveArticleRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// ArticleListResponse represents the response for listing articles
type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query  string   `json:"query"`
	Domain string   `json:"domain,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Read   *bool    `json:"read,omitempty"`
}
