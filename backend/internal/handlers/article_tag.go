package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"hotel-backend/internal/database"
)

type ArticleTagHandler struct{ DB *database.DB }

type articleTagResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	PostCount int    `json:"post_count"`
	IsSystem  bool   `json:"is_system"`
}

func NewArticleTagHandler(db *database.DB) *ArticleTagHandler { return &ArticleTagHandler{DB: db} }

func (h *ArticleTagHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Pool.Query(context.Background(), `
		SELECT c.id, c.name, c.sort_order, COUNT(pat.post_id)::int,
		       COALESCE(c.external_code = 'latest_posts', FALSE)
		FROM categories c LEFT JOIN post_article_tags pat ON pat.article_tag_id = c.id
		WHERE c.type = 'article_tag' AND c.parent_id IS NULL
		GROUP BY c.id, c.name, c.sort_order, c.external_code
		ORDER BY c.sort_order, c.id`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("failed to query article tags: %v", err)})
	}
	defer rows.Close()
	tags := []articleTagResponse{}
	for rows.Next() {
		var tag articleTagResponse
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.SortOrder, &tag.PostCount, &tag.IsSystem); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan article tag"})
		}
		tags = append(tags, tag)
	}
	return c.JSON(tags)
}

func (h *ArticleTagHandler) Get(c *fiber.Ctx) error {
	var tag articleTagResponse
	err := h.DB.Pool.QueryRow(context.Background(), `
		SELECT c.id, c.name, c.sort_order, COUNT(pat.post_id)::int,
		       COALESCE(c.external_code = 'latest_posts', FALSE)
		FROM categories c LEFT JOIN post_article_tags pat ON pat.article_tag_id = c.id
		WHERE c.id = $1 AND c.type = 'article_tag' AND c.parent_id IS NULL
		GROUP BY c.id, c.name, c.sort_order, c.external_code`, c.Params("id")).Scan(&tag.ID, &tag.Name, &tag.SortOrder, &tag.PostCount, &tag.IsSystem)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(404).JSON(fiber.Map{"error": "article tag not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to query article tag"})
	}
	return c.JSON(tag)
}

func cleanArticleTagInput(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > 30 {
		return "", fmt.Errorf("標籤名稱必須為 1 至 30 個字")
	}
	return name, nil
}

func articleTagSaveError(c *fiber.Ctx, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return c.Status(409).JSON(fiber.Map{"error": "標籤名稱已存在"})
	}
	return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("failed to save article tag: %v", err)})
}

func (h *ArticleTagHandler) Create(c *fiber.Ctx) error {
	var input struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	name, err := cleanArticleTagInput(input.Name)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	var tag articleTagResponse
	err = h.DB.Pool.QueryRow(context.Background(), `INSERT INTO categories (type,name,sort_order) VALUES ('article_tag',$1,$2) RETURNING id,name,sort_order`, name, input.SortOrder).Scan(&tag.ID, &tag.Name, &tag.SortOrder)
	if err != nil {
		return articleTagSaveError(c, err)
	}
	return c.Status(201).JSON(tag)
}

func (h *ArticleTagHandler) Update(c *fiber.Ctx) error {
	var system bool
	err := h.DB.Pool.QueryRow(context.Background(), `SELECT COALESCE(external_code = 'latest_posts', FALSE) FROM categories WHERE id=$1 AND type='article_tag'`, c.Params("id")).Scan(&system)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(404).JSON(fiber.Map{"error": "article tag not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to query article tag"})
	}
	if system {
		return c.Status(409).JSON(fiber.Map{"error": "最新文章為系統固定標籤，不可編輯"})
	}
	var input struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	name, err := cleanArticleTagInput(input.Name)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = h.DB.Pool.Exec(context.Background(), `UPDATE categories SET name=$1,sort_order=$2 WHERE id=$3 AND type='article_tag'`, name, input.SortOrder, c.Params("id"))
	if err != nil {
		return articleTagSaveError(c, err)
	}
	return c.JSON(fiber.Map{"message": "article tag updated"})
}

func (h *ArticleTagHandler) Usage(c *fiber.Ctx) error {
	rows, err := h.DB.Pool.Query(context.Background(), `SELECT p.id,p.title FROM post_article_tags pat JOIN posts p ON p.id=pat.post_id JOIN categories c ON c.id=pat.article_tag_id AND c.type='article_tag' WHERE pat.article_tag_id=$1 ORDER BY p.created_at DESC,p.id DESC`, c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to query article tag usage"})
	}
	defer rows.Close()
	posts := []fiber.Map{}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to scan article tag usage"})
		}
		posts = append(posts, fiber.Map{"id": id, "title": title})
	}
	return c.JSON(fiber.Map{"posts": posts, "total": len(posts)})
}

func (h *ArticleTagHandler) Delete(c *fiber.Ctx) error {
	if c.Query("confirm") != "true" {
		return c.Status(400).JSON(fiber.Map{"error": "deleting an article tag requires confirm=true"})
	}
	ctx := context.Background()
	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to begin transaction"})
	}
	defer tx.Rollback(ctx)
	var system bool
	err = tx.QueryRow(ctx, `SELECT COALESCE(external_code = 'latest_posts', FALSE) FROM categories WHERE id=$1 AND type='article_tag'`, c.Params("id")).Scan(&system)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(404).JSON(fiber.Map{"error": "article tag not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to validate article tag"})
	}
	if system {
		return c.Status(409).JSON(fiber.Map{"error": "最新文章為系統固定標籤，不可刪除"})
	}
	result, err := tx.Exec(ctx, `DELETE FROM post_article_tags WHERE article_tag_id=$1`, c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to detach article tag"})
	}
	if _, err = tx.Exec(ctx, `DELETE FROM categories WHERE id=$1 AND type='article_tag'`, c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete article tag"})
	}
	if err = tx.Commit(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to commit article tag deletion"})
	}
	return c.JSON(fiber.Map{"message": "article tag deleted", "detached_posts": result.RowsAffected()})
}
