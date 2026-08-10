package models

import (
	"time"
)

// Post represents a blog post database record and JSON structure.
type Post struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	Tags           []string     `json:"tags"` // Maps JSONB array in PostgreSQL
	Image          string       `json:"image"`
	Content        string       `json:"content"`
	AdLink         string       `json:"ad_link"`
	SEOTitle       string       `json:"seo_title"`
	SEOKeywords    string       `json:"seo_keywords"`
	SEODescription string       `json:"seo_description"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	ArticleTags    []ArticleTag `json:"article_tags"`
	ArticleTagIDs  []int        `json:"article_tag_ids"`
}

type ArticleTag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsSystem bool   `json:"is_system"`
}
