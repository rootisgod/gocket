package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"gocket/internal/models"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// CreateArticle saves a new article to the database
func (r *ArticleRepository) CreateArticle(article *models.Article) error {
	tagsJSON, _ := json.Marshal(article.Tags)

	query := `
		INSERT INTO articles (url, title, content, excerpt, author, published_date, 
		                     saved_date, read_status, tags, word_count, reading_time, 
		                     thumbnail_url, domain, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		article.URL,
		article.Title,
		article.Content,
		article.Excerpt,
		article.Author,
		article.PublishedDate,
		article.SavedDate,
		article.ReadStatus,
		tagsJSON,
		article.WordCount,
		article.ReadingTime,
		article.ThumbnailURL,
		article.Domain,
		article.CreatedAt,
		article.UpdatedAt,
	)

	return err
}

// GetArticleByID retrieves an article by its ID
func (r *ArticleRepository) GetArticleByID(id int) (*models.Article, error) {
	query := `SELECT id, url, title, content, excerpt, author, published_date, 
	                 saved_date, read_status, tags, word_count, reading_time, 
	                 thumbnail_url, domain, created_at, updated_at
	          FROM articles WHERE id = ?`

	row := r.db.QueryRow(query, id)

	article := &models.Article{}
	var tagsJSON string
	var publishedDate sql.NullTime

	err := row.Scan(
		&article.ID,
		&article.URL,
		&article.Title,
		&article.Content,
		&article.Excerpt,
		&article.Author,
		&publishedDate,
		&article.SavedDate,
		&article.ReadStatus,
		&tagsJSON,
		&article.WordCount,
		&article.ReadingTime,
		&article.ThumbnailURL,
		&article.Domain,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if publishedDate.Valid {
		article.PublishedDate = &publishedDate.Time
	}

	json.Unmarshal([]byte(tagsJSON), &article.Tags)

	return article, nil
}

// GetAllArticles retrieves all articles with pagination
func (r *ArticleRepository) GetAllArticles(page, perPage int) ([]models.Article, int, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM articles"
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get articles
	query := `SELECT id, url, title, content, excerpt, author, published_date, 
	                 saved_date, read_status, tags, word_count, reading_time, 
	                 thumbnail_url, domain, created_at, updated_at
	          FROM articles 
	          ORDER BY saved_date DESC 
	          LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		article := models.Article{}
		var tagsJSON string
		var publishedDate sql.NullTime

		err := rows.Scan(
			&article.ID,
			&article.URL,
			&article.Title,
			&article.Content,
			&article.Excerpt,
			&article.Author,
			&publishedDate,
			&article.SavedDate,
			&article.ReadStatus,
			&tagsJSON,
			&article.WordCount,
			&article.ReadingTime,
			&article.ThumbnailURL,
			&article.Domain,
			&article.CreatedAt,
			&article.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		if publishedDate.Valid {
			article.PublishedDate = &publishedDate.Time
		}

		json.Unmarshal([]byte(tagsJSON), &article.Tags)
		articles = append(articles, article)
	}

	return articles, total, nil
}

// UpdateArticle updates an existing article
func (r *ArticleRepository) UpdateArticle(article *models.Article) error {
	tagsJSON, _ := json.Marshal(article.Tags)

	query := `
		UPDATE articles 
		SET title = ?, content = ?, excerpt = ?, author = ?, published_date = ?, 
		    read_status = ?, tags = ?, word_count = ?, reading_time = ?, 
		    thumbnail_url = ?, domain = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		article.Title,
		article.Content,
		article.Excerpt,
		article.Author,
		article.PublishedDate,
		article.ReadStatus,
		tagsJSON,
		article.WordCount,
		article.ReadingTime,
		article.ThumbnailURL,
		article.Domain,
		time.Now(),
		article.ID,
	)

	return err
}

// DeleteArticle removes an article from the database
func (r *ArticleRepository) DeleteArticle(id int) error {
	query := "DELETE FROM articles WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
