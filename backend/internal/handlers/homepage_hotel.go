package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type HomepageHotelHandler struct {
	DB *database.DB
}

func NewHomepageHotelHandler(db *database.DB) *HomepageHotelHandler {
	return &HomepageHotelHandler{DB: db}
}

type HomepageHotelResponseItem struct {
	City      string   `json:"city"`
	SortOrder int      `json:"sort_order"`
	HotelID   string   `json:"hotel_id"`
	Name      string   `json:"name"`
	Area      string   `json:"area"`
	Images    []string `json:"images"`
}

// GetHomepageHotels retrieves all custom featured homepage hotels.
// GET /api/homepage-hotels
func (h *HomepageHotelHandler) GetHomepageHotels(c *fiber.Ctx) error {
	ctx := context.Background()

	query := `
		SELECT hh.city, hh.sort_order, h.id, h.name, COALESCE(h.area, ''),
		       COALESCE((
		           SELECT jsonb_agg(hi.url ORDER BY hi.sort_order, hi.id)
		           FROM hotel_images hi WHERE hi.hotel_id = h.id
		       ), '[]'::jsonb)
		FROM homepage_hotels hh
		JOIN hotels h ON h.id = hh.hotel_id
		ORDER BY hh.city, hh.sort_order`

	rows, err := h.DB.Pool.Query(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query homepage hotels: %v", err)})
	}
	defer rows.Close()

	items := []HomepageHotelResponseItem{}
	for rows.Next() {
		var item HomepageHotelResponseItem
		var imagesJSON []byte
		err := rows.Scan(&item.City, &item.SortOrder, &item.HotelID, &item.Name, &item.Area, &imagesJSON)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan item: %v", err)})
		}
		if len(imagesJSON) > 0 {
			_ = json.Unmarshal(imagesJSON, &item.Images)
		}
		if item.Images == nil {
			item.Images = []string{}
		}
		items = append(items, item)
	}

	return c.JSON(items)
}

type SaveHomepageHotelsReq struct {
	Selections []models.HomepageHotel `json:"selections"`
}

// UpdateHomepageHotels bulk updates homepage hotel slots.
// PUT /api/homepage-hotels
func (h *HomepageHotelHandler) UpdateHomepageHotels(c *fiber.Ctx) error {
	var req SaveHomepageHotelsReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()
	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to start transaction: %v", err)})
	}
	defer tx.Rollback(ctx)

	// Clear existing selections
	_, err = tx.Exec(ctx, "DELETE FROM homepage_hotels")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to clear selections: %v", err)})
	}

	// Insert new selections
	insertQuery := `
		INSERT INTO homepage_hotels (city, sort_order, hotel_id)
		VALUES ($1, $2, $3)`

	for _, sel := range req.Selections {
		if sel.City == "" || sel.HotelID == "" || sel.SortOrder < 1 || sel.SortOrder > 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid selection values"})
		}
		_, err = tx.Exec(ctx, insertQuery, sel.City, sel.SortOrder, sel.HotelID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to insert selection: %v", err)})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to commit transaction: %v", err)})
	}

	return c.JSON(fiber.Map{"success": true})
}
