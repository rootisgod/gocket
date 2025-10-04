package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gocket/internal/models"
	"gocket/internal/scraper/extractors"

	"github.com/PuerkitoBio/goquery"
)

type ContentProcessor struct {
	client *http.Client
}

func NewContentProcessor() *ContentProcessor {
	return &ContentProcessor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessURL fetches and processes a web page URL
func (cp *ContentProcessor) ProcessURL(urlStr string) (*models.Article, error) {
	// 1. Fetch HTML content
	html, err := cp.fetchHTML(urlStr)
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
	article := &models.Article{
		URL:           urlStr,
		Title:         metadata.Title,
		Content:       content,
		Excerpt:       cp.generateExcerpt(content),
		Author:        metadata.Author,
		PublishedDate: metadata.PublishedDate,
		WordCount:     wordCount,
		ReadingTime:   readingTime,
		Domain:        cp.extractDomain(urlStr),
		ReadStatus:    false,
		Tags:          []string{},
	}

	return article, nil
}

// fetchHTML downloads the HTML content from a URL
func (cp *ContentProcessor) fetchHTML(urlStr string) (string, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Gocket/1.0")

	resp, err := cp.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body := make([]byte, 0, 1024*1024) // 1MB buffer
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(body), nil
}

// extractContent parses HTML and extracts clean content
func (cp *ContentProcessor) extractContent(html string) (string, *extractors.Metadata, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", nil, err
	}

	// Extract metadata
	metadata := extractors.ExtractMetadata(doc)

	// Extract main content
	content := extractors.ExtractMainContent(doc)

	return content, metadata, nil
}

// calculateWordCount estimates word count from content
func (cp *ContentProcessor) calculateWordCount(content string) int {
	words := strings.Fields(content)
	return len(words)
}

// estimateReadingTime calculates reading time in minutes
func (cp *ContentProcessor) estimateReadingTime(wordCount int) int {
	// Average reading speed: 200 words per minute
	minutes := wordCount / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

// generateExcerpt creates a short excerpt from content
func (cp *ContentProcessor) generateExcerpt(content string) string {
	words := strings.Fields(content)
	if len(words) <= 50 {
		return content
	}
	return strings.Join(words[:50], " ") + "..."
}

// extractDomain extracts domain from URL
func (cp *ContentProcessor) extractDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return u.Host
}
