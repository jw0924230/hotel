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

type HotelTagHandler struct {
	DB *database.DB
}

type hotelTagResponse struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	SortOrder         int    `json:"sort_order"`
	HotelCount        int    `json:"hotel_count"`
	EnabledHotelCount int    `json:"enabled_hotel_count"`
}

func NewHotelTagHandler(db *database.DB) *HotelTagHandler {
	return &HotelTagHandler{DB: db}
}

func (h *HotelTagHandler) List(c *fiber.Ctx) error {
	rows, err := h.DB.Pool.Query(context.Background(), `
		SELECT c.id, c.name, c.sort_order,
		       COUNT(ht.hotel_id)::int,
		       COUNT(ht.hotel_id) FILTER (WHERE h.is_disabled = FALSE)::int
		FROM categories c
		LEFT JOIN hotel_tags ht ON ht.tag_id = c.id
		LEFT JOIN hotels h ON h.id = ht.hotel_id
		WHERE c.type = 'hotel_tag' AND c.parent_id IS NULL
		GROUP BY c.id, c.name, c.sort_order
		ORDER BY c.sort_order ASC, c.id ASC`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotel tags: %v", err)})
	}
	defer rows.Close()

	tags := []hotelTagResponse{}
	for rows.Next() {
		var tag hotelTagResponse
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.SortOrder, &tag.HotelCount, &tag.EnabledHotelCount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan hotel tag: %v", err)})
		}
		tags = append(tags, tag)
	}
	return c.JSON(tags)
}

func (h *HotelTagHandler) Get(c *fiber.Ctx) error {
	var tag hotelTagResponse
	err := h.DB.Pool.QueryRow(context.Background(), `
		SELECT c.id, c.name, c.sort_order,
		       COUNT(ht.hotel_id)::int,
		       COUNT(ht.hotel_id) FILTER (WHERE h.is_disabled = FALSE)::int
		FROM categories c
		LEFT JOIN hotel_tags ht ON ht.tag_id = c.id
		LEFT JOIN hotels h ON h.id = ht.hotel_id
		WHERE c.id = $1 AND c.type = 'hotel_tag' AND c.parent_id IS NULL
		GROUP BY c.id, c.name, c.sort_order`, c.Params("id")).Scan(
		&tag.ID, &tag.Name, &tag.SortOrder, &tag.HotelCount, &tag.EnabledHotelCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "hotel tag not found"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotel tag: %v", err)})
	}
	return c.JSON(tag)
}

func cleanHotelTagInput(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > 30 {
		return "", fmt.Errorf("標籤名稱必須為 1 至 30 個字")
	}
	return name, nil
}

func hotelTagConflict(c *fiber.Ctx, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "標籤名稱已存在"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to save hotel tag: %v", err)})
}

func (h *HotelTagHandler) Create(c *fiber.Ctx) error {
	var input struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	name, err := cleanHotelTagInput(input.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var tag hotelTagResponse
	err = h.DB.Pool.QueryRow(context.Background(), `
		INSERT INTO categories (type, name, sort_order)
		VALUES ('hotel_tag', $1, $2)
		RETURNING id, name, sort_order`, name, input.SortOrder,
	).Scan(&tag.ID, &tag.Name, &tag.SortOrder)
	if err != nil {
		return hotelTagConflict(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *HotelTagHandler) Update(c *fiber.Ctx) error {
	var input struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	name, err := cleanHotelTagInput(input.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.DB.Pool.Exec(context.Background(), `
		UPDATE categories SET name = $1, sort_order = $2
		WHERE id = $3 AND type = 'hotel_tag'`, name, input.SortOrder, c.Params("id"),
	)
	if err != nil {
		return hotelTagConflict(c, err)
	}
	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "hotel tag not found"})
	}
	return c.JSON(fiber.Map{"message": "hotel tag updated"})
}

func (h *HotelTagHandler) Usage(c *fiber.Ctx) error {
	rows, err := h.DB.Pool.Query(context.Background(), `
		SELECT h.id, h.name, h.is_disabled
		FROM hotel_tags ht
		JOIN hotels h ON h.id = ht.hotel_id
		JOIN categories c ON c.id = ht.tag_id AND c.type = 'hotel_tag'
		WHERE ht.tag_id = $1
		ORDER BY h.name ASC, h.id ASC`, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotel tag usage: %v", err)})
	}
	defer rows.Close()

	hotels := []fiber.Map{}
	for rows.Next() {
		var id, name string
		var disabled bool
		if err := rows.Scan(&id, &name, &disabled); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan hotel tag usage: %v", err)})
		}
		hotels = append(hotels, fiber.Map{"id": id, "name": name, "is_disabled": disabled})
	}
	return c.JSON(fiber.Map{"hotels": hotels, "total": len(hotels)})
}

func (h *HotelTagHandler) Delete(c *fiber.Ctx) error {
	if c.Query("confirm") != "true" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "deleting a hotel tag requires confirm=true"})
	}

	ctx := context.Background()
	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to begin transaction"})
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1 AND type = 'hotel_tag')`, c.Params("id")).Scan(&exists); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to validate hotel tag"})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "hotel tag not found"})
	}

	var detached int
	result, err := tx.Exec(ctx, `DELETE FROM hotel_tags WHERE tag_id = $1`, c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to detach hotel tag"})
	}
	detached = int(result.RowsAffected())
	if _, err = tx.Exec(ctx, `DELETE FROM categories WHERE id = $1 AND type = 'hotel_tag'`, c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete hotel tag"})
	}
	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to commit hotel tag deletion"})
	}
	return c.JSON(fiber.Map{"message": "hotel tag deleted", "detached_hotels": detached})
}
