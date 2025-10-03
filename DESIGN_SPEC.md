# Gocket - Design and Implementation Specification

## Overview

Gocket is a local Pocket clone that allows users to save web pages for offline reading. It's designed to run as a single Go application with a web interface, using SQLite for data persistence and designed for containerized deployment.

## Core Requirements

- **Language**: Go
- **Web Interface**: Local web server for browsing saved articles
- **URL Processing**: Accept web page URLs and save snapshots
- **Data Storage**: SQLite database for article persistence
- **Deployment**: Single Docker container or Kubernetes pod
- **Scalability**: Designed to support multiple users when hosted

## System Architecture

### High-Level Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   Web Server    │    │   SQLite DB     │
│   (Browser)     │◄──►│   (Go HTTP)     │◄──►│   (Articles)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Web Scraper    │
                       │  (goquery)      │
                       └─────────────────┘
```

### Component Responsibilities

1. **Web Server**: HTTP server handling API requests and serving static content
2. **Article Processor**: Fetches and processes web pages for storage
3. **Database Layer**: SQLite operations for article CRUD
4. **Web Interface**: HTML/CSS/JS frontend for article browsing

## Database Schema

### Articles Table

```sql
CREATE TABLE articles (
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
);

CREATE INDEX idx_articles_url ON articles(url);
CREATE INDEX idx_articles_saved_date ON articles(saved_date);
CREATE INDEX idx_articles_read_status ON articles(read_status);
CREATE INDEX idx_articles_domain ON articles(domain);
```

### Tags Table (Optional - for normalized tag storage)

```sql
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE article_tags (
    article_id INTEGER,
    tag_id INTEGER,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
```

## API Endpoints

### Article Management

```
POST   /api/articles              # Save new article from URL
GET    /api/articles              # List all articles (with pagination)
GET    /api/articles/:id          # Get specific article
PUT    /api/articles/:id          # Update article (mark as read, add tags)
DELETE /api/articles/:id          # Delete article
GET    /api/articles/search       # Search articles by title/content
```

### Article Processing

```
POST   /api/articles/fetch        # Fetch and process URL
GET    /api/articles/:id/content  # Get article content for reading
```

### Statistics and Management

```
GET    /api/stats                 # Get reading statistics
GET    /api/health                # Health check endpoint
```

## Data Models

### Article Model

```go
type Article struct {
    ID            int       `json:"id" db:"id"`
    URL           string    `json:"url" db:"url"`
    Title         string    `json:"title" db:"title"`
    Content       string    `json:"content" db:"content"`
    Excerpt       string    `json:"excerpt" db:"excerpt"`
    Author        string    `json:"author" db:"author"`
    PublishedDate *time.Time `json:"published_date" db:"published_date"`
    SavedDate     time.Time `json:"saved_date" db:"saved_date"`
    ReadStatus    bool      `json:"read_status" db:"read_status"`
    Tags          []string  `json:"tags" db:"tags"`
    WordCount     int       `json:"word_count" db:"word_count"`
    ReadingTime   int       `json:"reading_time" db:"reading_time"`
    ThumbnailURL  string    `json:"thumbnail_url" db:"thumbnail_url"`
    Domain        string    `json:"domain" db:"domain"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
```

### API Request/Response Models

```go
type SaveArticleRequest struct {
    URL string `json:"url" validate:"required,url"`
}

type ArticleListResponse struct {
    Articles []Article `json:"articles"`
    Total    int       `json:"total"`
    Page     int       `json:"page"`
    PerPage  int       `json:"per_page"`
}

type SearchRequest struct {
    Query  string `json:"query"`
    Domain string `json:"domain,omitempty"`
    Tags   []string `json:"tags,omitempty"`
    Read   *bool   `json:"read,omitempty"`
}
```

## Web Scraping Implementation

### Content Extraction Strategy

1. **HTML Fetching**: Use `net/http` to fetch web pages
2. **Content Parsing**: Use `goquery` for HTML parsing
3. **Content Cleaning**: Remove ads, navigation, comments
4. **Metadata Extraction**: Extract title, author, published date
5. **Text Processing**: Calculate word count and reading time

### Content Processing Pipeline

```go
type ContentProcessor struct {
    client *http.Client
}

func (cp *ContentProcessor) ProcessURL(url string) (*Article, error) {
    // 1. Fetch HTML content
    html, err := cp.fetchHTML(url)
    if err != nil {
        return nil, err
    }
    
    // 2. Parse and clean content
    content, metadata, err := cp.extractContent(html)
    if err != nil {
        return nil, err
    }
    
    // 3. Calculate metrics
    wordCount := cp.calculateWordCount(content)
    readingTime := cp.estimateReadingTime(wordCount)
    
    // 4. Create article object
    article := &Article{
        URL:         url,
        Title:       metadata.Title,
        Content:     content,
        Author:      metadata.Author,
        PublishedDate: metadata.PublishedDate,
        WordCount:   wordCount,
        ReadingTime: readingTime,
        Domain:      cp.extractDomain(url),
    }
    
    return article, nil
}
```

## Web Interface Design

### Frontend Technology Stack

- **HTML5**: Semantic markup
- **CSS3**: Modern styling with CSS Grid/Flexbox
- **Vanilla JavaScript**: No external dependencies
- **Progressive Web App**: Offline capability

### Key UI Components

1. **Article List View**: Grid/list toggle, search, filters
2. **Article Reader**: Clean reading interface
3. **Add Article**: URL input with preview
4. **Settings**: Reading preferences, export options

### Responsive Design

- Mobile-first approach
- Touch-friendly interface
- Dark/light theme support
- Readable typography

## Project Structure

```
gocket/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   ├── middleware/          # HTTP middleware
│   │   └── routes/              # Route definitions
│   ├── database/
│   │   ├── migrations/          # SQL migration files
│   │   └── repository/          # Database operations
│   ├── scraper/
│   │   ├── processor.go         # Content processing
│   │   └── extractors.go        # Content extraction
│   └── models/
│       └── article.go           # Data models
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/
│       └── index.html
├── docker/
│   └── Dockerfile
├── k8s/
│   └── deployment.yaml
├── go.mod
├── go.sum
└── README.md
```

## Configuration

### Environment Variables

```bash
# Server Configuration
PORT=8080
HOST=0.0.0.0

# Database Configuration
DB_PATH=/data/gocket.db

# Scraping Configuration
USER_AGENT="Gocket/1.0"
REQUEST_TIMEOUT=30s
MAX_CONTENT_SIZE=10MB

# Security
CORS_ORIGINS=*
RATE_LIMIT=100/hour
```

## Deployment Strategy

### Docker Configuration

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gocket cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gocket .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./gocket"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gocket
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gocket
  template:
    metadata:
      labels:
        app: gocket
    spec:
      containers:
      - name: gocket
        image: gocket:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DB_PATH
          value: "/data/gocket.db"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: gocket-data
```

## Performance Considerations

### Database Optimization

- Connection pooling
- Prepared statements
- Proper indexing
- Query optimization

### Caching Strategy

- In-memory article cache
- HTTP response caching
- Static asset caching

### Scalability Features

- Horizontal scaling with load balancer
- Database connection pooling
- Rate limiting
- Content compression

## Security Considerations

### Input Validation

- URL validation and sanitization
- SQL injection prevention
- XSS protection
- CSRF protection

### Network Security

- HTTPS enforcement
- CORS configuration
- Rate limiting
- Request size limits

## Testing Strategy

### Unit Tests

- Database operations
- Content processing
- API handlers
- Utility functions

### Integration Tests

- End-to-end API testing
- Database integration
- Web scraping functionality

### Performance Tests

- Load testing
- Database performance
- Memory usage profiling

## Monitoring and Logging

### Logging

- Structured logging with levels
- Request/response logging
- Error tracking
- Performance metrics

### Health Checks

- Database connectivity
- External service availability
- Memory usage
- Disk space

## Future Enhancements

### Phase 2 Features

- User authentication
- Article sharing
- Export functionality (PDF, EPUB)
- Mobile app
- Browser extension

### Phase 3 Features

- Multi-user support
- Article recommendations
- Reading statistics
- Social features
- API for third-party integrations

## Implementation Timeline

### Phase 1 (MVP) - 2-3 weeks
- Basic web server
- SQLite database setup
- Article saving functionality
- Simple web interface
- Docker deployment

### Phase 2 (Enhanced) - 2-3 weeks
- Advanced content processing
- Search and filtering
- Improved UI/UX
- Performance optimizations
- Kubernetes deployment

### Phase 3 (Production) - 2-3 weeks
- Security hardening
- Monitoring and logging
- Testing coverage
- Documentation
- Performance tuning

## Dependencies

### Core Dependencies

```go
// Web framework
github.com/gorilla/mux

// Database
github.com/mattn/go-sqlite3

// HTML parsing
github.com/PuerkitoBio/goquery

// Configuration
github.com/spf13/viper

// Validation
github.com/go-playground/validator

// Logging
github.com/sirupsen/logrus
```

### Development Dependencies

```go
// Testing
github.com/stretchr/testify

// Code generation
github.com/golang-migrate/migrate

// Linting
golang.org/x/lint
```

This specification provides a comprehensive foundation for implementing Gocket as a local Pocket clone with modern Go practices and containerized deployment capabilities.
