package extractors

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Metadata contains extracted article metadata
type Metadata struct {
	Title         string
	Author        string
	PublishedDate *time.Time
}

// ExtractMetadata extracts title, author, and published date from document
func ExtractMetadata(doc *goquery.Document) *Metadata {
	metadata := &Metadata{}

	// Extract title
	title := doc.Find("title").First().Text()
	if title == "" {
		title = doc.Find("h1").First().Text()
	}
	metadata.Title = strings.TrimSpace(title)

	// Extract author
	author := doc.Find("meta[name='author']").AttrOr("content", "")
	if author == "" {
		author = doc.Find("[rel='author']").Text()
	}
	if author == "" {
		author = doc.Find(".author").First().Text()
	}
	metadata.Author = strings.TrimSpace(author)

	// Extract published date
	publishedDate := doc.Find("meta[property='article:published_time']").AttrOr("content", "")
	if publishedDate == "" {
		publishedDate = doc.Find("meta[name='date']").AttrOr("content", "")
	}
	if publishedDate == "" {
		publishedDate = doc.Find("time[datetime]").AttrOr("datetime", "")
	}

	if publishedDate != "" {
		if t, err := time.Parse(time.RFC3339, publishedDate); err == nil {
			metadata.PublishedDate = &t
		}
	}

	return metadata
}

// ExtractMainContent extracts the main article content
func ExtractMainContent(doc *goquery.Document) string {
	// Remove unwanted elements
	doc.Find("script, style, nav, header, footer, aside, .advertisement, .ads, .social-share").Remove()

	// Try to find main content area
	var content string

	// Look for common article selectors
	selectors := []string{
		"article",
		".article-content",
		".post-content",
		".entry-content",
		".content",
		"main",
		"#content",
		".main-content",
	}

	for _, selector := range selectors {
		if element := doc.Find(selector).First(); element.Length() > 0 {
			content = element.Text()
			break
		}
	}

	// Fallback to body if no specific content found
	if content == "" {
		content = doc.Find("body").Text()
	}

	// Clean up whitespace
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.ReplaceAll(content, "\t", " ")

	// Remove multiple spaces
	for strings.Contains(content, "  ") {
		content = strings.ReplaceAll(content, "  ", " ")
	}

	return strings.TrimSpace(content)
}
