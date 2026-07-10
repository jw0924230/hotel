package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"
)

// PostHandler holds database dependencies for post routes.
type PostHandler struct {
	DB *database.DB
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(db *database.DB) *PostHandler {
	return &PostHandler{DB: db}
}

// PostUpsertRequest defines the request body for creating or updating a post.
type PostUpsertRequest struct {
	Title          string   `json:"title"`
	Tags           []string `json:"tags"`
	Image          string   `json:"image"`
	Content        string   `json:"content"`
	AdLink         string   `json:"ad_link"`
	SEOTitle       string   `json:"seo_title"`
	SEOKeywords    string   `json:"seo_keywords"`
	SEODescription string   `json:"seo_description"`
}

// List retrieves a list of posts with pagination, search, and tag filters.
// GET /api/posts
func (h *PostHandler) List(c *fiber.Ctx) error {
	ctx := context.Background()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := strings.TrimSpace(c.Query("search", ""))
	tag := strings.TrimSpace(c.Query("tag", ""))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Build total count query
	countQuery := `SELECT COUNT(*) FROM posts WHERE 1=1`
	var countArgs []interface{}
	argCount := 1

	if search != "" {
		countQuery += fmt.Sprintf(" AND title ILIKE $%d", argCount)
		countArgs = append(countArgs, "%"+search+"%")
		argCount++
	}
	if tag != "" {
		tagJSON, _ := json.Marshal([]string{tag})
		countQuery += fmt.Sprintf(" AND tags @> $%d::jsonb", argCount)
		countArgs = append(countArgs, string(tagJSON))
		argCount++
	}

	var total int
	err := h.DB.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to count posts: %v", err)})
	}

	// Build select query
	selectQuery := `
		SELECT id, title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description, created_at, updated_at
		FROM posts
		WHERE 1=1`
	var selectArgs []interface{}
	argSelectCount := 1

	if search != "" {
		selectQuery += fmt.Sprintf(" AND title ILIKE $%d", argSelectCount)
		selectArgs = append(selectArgs, "%"+search+"%")
		argSelectCount++
	}
	if tag != "" {
		tagJSON, _ := json.Marshal([]string{tag})
		selectQuery += fmt.Sprintf(" AND tags @> $%d::jsonb", argSelectCount)
		selectArgs = append(selectArgs, string(tagJSON))
		argSelectCount++
	}

	selectQuery += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", argSelectCount, argSelectCount+1)
	selectArgs = append(selectArgs, limit, offset)

	rows, err := h.DB.Pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query posts: %v", err)})
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		var tagsJSON []byte
		err := rows.Scan(
			&post.ID, &post.Title, &tagsJSON, &post.Image, &post.Content, &post.AdLink,
			&post.SEOTitle, &post.SEOKeywords, &post.SEODescription, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan post: %v", err)})
		}

		if len(tagsJSON) > 0 {
			_ = json.Unmarshal(tagsJSON, &post.Tags)
		}
		if post.Tags == nil {
			post.Tags = []string{}
		}

		posts = append(posts, post)
	}

	return c.JSON(fiber.Map{
		"data":  posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Get retrieves a single post by ID.
// GET /api/posts/:id
func (h *PostHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
	}

	ctx := context.Background()
	var post models.Post
	var tagsJSON []byte

	query := `
		SELECT id, title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description, created_at, updated_at
		FROM posts WHERE id = $1`

	err = h.DB.Pool.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.Title, &tagsJSON, &post.Image, &post.Content, &post.AdLink,
		&post.SEOTitle, &post.SEOKeywords, &post.SEODescription, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query post: %v", err)})
	}

	if len(tagsJSON) > 0 {
		_ = json.Unmarshal(tagsJSON, &post.Tags)
	}
	if post.Tags == nil {
		post.Tags = []string{}
	}

	return c.JSON(post)
}

// Create creates a new blog post.
// POST /api/posts
func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req PostUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文章標題不能為空"})
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文章內容不能為空"})
	}

	// Validate tags limit (max 3 tags)
	cleanTags := []string{}
	for _, t := range req.Tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleanTags = append(cleanTags, trimmed)
		}
	}
	if len(cleanTags) > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "分類標籤最多僅能設定三個"})
	}

	tagsJSON, err := json.Marshal(cleanTags)
	if err != nil {
		tagsJSON = []byte("[]")
	}

	ctx := context.Background()
	var post models.Post
	query := `
		INSERT INTO posts (
			title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description, created_at, updated_at`

	var returnedTagsJSON []byte
	err = h.DB.Pool.QueryRow(ctx, query,
		req.Title, string(tagsJSON), req.Image, req.Content, req.AdLink,
		req.SEOTitle, req.SEOKeywords, req.SEODescription,
	).Scan(
		&post.ID, &post.Title, &returnedTagsJSON, &post.Image, &post.Content, &post.AdLink,
		&post.SEOTitle, &post.SEOKeywords, &post.SEODescription, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to create post: %v", err)})
	}

	if len(returnedTagsJSON) > 0 {
		_ = json.Unmarshal(returnedTagsJSON, &post.Tags)
	}

	return c.JSON(post)
}

// Update updates an existing blog post.
// PUT /api/posts/:id
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
	}

	var req PostUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文章標題不能為空"})
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "文章內容不能為空"})
	}

	// Validate tags limit (max 3 tags)
	cleanTags := []string{}
	for _, t := range req.Tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleanTags = append(cleanTags, trimmed)
		}
	}
	if len(cleanTags) > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "分類標籤最多僅能設定三個"})
	}

	tagsJSON, err := json.Marshal(cleanTags)
	if err != nil {
		tagsJSON = []byte("[]")
	}

	ctx := context.Background()
	query := `
		UPDATE posts SET
			title = $1,
			tags = $2,
			image = $3,
			content = $4,
			ad_link = $5,
			seo_title = $6,
			seo_keywords = $7,
			seo_description = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9`

	_, err = h.DB.Pool.Exec(ctx, query,
		req.Title, string(tagsJSON), req.Image, req.Content, req.AdLink,
		req.SEOTitle, req.SEOKeywords, req.SEODescription, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to update post: %v", err)})
	}

	// Retrieve updated object
	var post models.Post
	var returnedTagsJSON []byte
	queryGet := `
		SELECT id, title, tags, image, content, ad_link, seo_title, seo_keywords, seo_description, created_at, updated_at
		FROM posts WHERE id = $1`

	err = h.DB.Pool.QueryRow(ctx, queryGet, id).Scan(
		&post.ID, &post.Title, &returnedTagsJSON, &post.Image, &post.Content, &post.AdLink,
		&post.SEOTitle, &post.SEOKeywords, &post.SEODescription, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query updated post: %v", err)})
	}

	if len(returnedTagsJSON) > 0 {
		_ = json.Unmarshal(returnedTagsJSON, &post.Tags)
	}

	return c.JSON(post)
}

// Delete deletes a blog post.
// DELETE /api/posts/:id
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid post ID"})
	}

	ctx := context.Background()
	query := `DELETE FROM posts WHERE id = $1`
	_, err = h.DB.Pool.Exec(ctx, query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to delete post: %v", err)})
	}

	return c.JSON(fiber.Map{"message": "文章已成功刪除"})
}
