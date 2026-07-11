package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"
)

// CategoryHandler holds database dependencies for category routes.
type CategoryHandler struct {
	DB *database.DB
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(db *database.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// List retrieves categories filtered by type from the database.
// GET /api/categories?type=city
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	categoryType := c.Query("type", "")

	ctx := context.Background()

	query := `SELECT id, type, name, sort_order, created_at FROM categories`
	var args []interface{}

	if categoryType != "" {
		query += ` WHERE type = $1`
		args = append(args, categoryType)
	}

	query += ` ORDER BY sort_order ASC, id ASC`

	rows, err := h.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query categories: %v", err)})
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Type, &cat.Name, &cat.SortOrder, &cat.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan category: %v", err)})
		}
		categories = append(categories, cat)
	}

	return c.JSON(categories)
}

// Create adds a new category.
// POST /api/categories
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var input struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Type == "" || input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "type and name are required"})
	}

	ctx := context.Background()
	var cat models.Category

	err := h.DB.Pool.QueryRow(ctx,
		`INSERT INTO categories (type, name, sort_order) VALUES ($1, $2, $3) RETURNING id, type, name, sort_order, created_at`,
		input.Type, input.Name, input.SortOrder,
	).Scan(&cat.ID, &cat.Type, &cat.Name, &cat.SortOrder, &cat.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to create category: %v", err)})
	}

	return c.Status(fiber.StatusCreated).JSON(cat)
}

// Update modifies an existing category.
// PUT /api/categories/:id
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	var input struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	ctx := context.Background()
	result, err := h.DB.Pool.Exec(ctx,
		`UPDATE categories SET name = $1, sort_order = $2 WHERE id = $3`,
		input.Name, input.SortOrder, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to update category: %v", err)})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "category not found"})
	}

	return c.JSON(fiber.Map{"message": "category updated"})
}

// Delete removes a category.
// DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	ctx := context.Background()
	result, err := h.DB.Pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to delete category: %v", err)})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "category not found"})
	}

	return c.JSON(fiber.Map{"message": "category deleted"})
}

// RegionResponse represents the structure of a region grouping in the JSON file.
type RegionResponse struct {
	Name   string   `json:"name"`
	Cities []string `json:"cities"`
}

// StaticRegions hardcodes regions.json content as a static backend parameter.
var StaticRegions = []RegionResponse{
	{"北部", []string{"基隆", "台北", "新北", "桃園", "新竹", "宜蘭"}},
	{"中部", []string{"苗栗", "台中", "彰化", "南投", "雲林"}},
	{"南部", []string{"嘉義", "台南", "高雄", "屏東"}},
	{"東部", []string{"花蓮", "台東"}},
	{"海外 / 外島", []string{"澎湖", "金門", "馬祖", "其他"}},
}

// Regions serves the static regions definition directly from memory.
// GET /api/regions
func (h *CategoryHandler) Regions(c *fiber.Ctx) error {
	return c.JSON(StaticRegions)
}

type CombinedCityInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CombinedRegionInfo struct {
	Name   string             `json:"name"`
	Cities []CombinedCityInfo `json:"cities"`
}

type CombinedLocationsResponse struct {
	Cities  []CombinedCityInfo   `json:"cities"`
	Regions []CombinedRegionInfo `json:"regions"`
}

// CombinedRegions combines categories from DB and static regions to serve a single locations API.
// GET /api/regions/combined
func (h *CategoryHandler) CombinedRegions(c *fiber.Ctx) error {
	ctx := context.Background()

	// 1. Fetch categories (cities) from database
	query := `SELECT id, name, sort_order FROM categories WHERE type = 'city' ORDER BY sort_order ASC, id ASC`
	rows, err := h.DB.Pool.Query(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query categories: %v", err)})
	}
	defer rows.Close()

	var cities []CombinedCityInfo
	cityMap := make(map[string]CombinedCityInfo)

	for rows.Next() {
		var id int
		var name string
		var sortOrder int
		if err := rows.Scan(&id, &name, &sortOrder); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan category: %v", err)})
		}

		// Use sort_order if it exists/non-zero, otherwise fallback to id
		finalID := sortOrder
		if finalID == 0 {
			finalID = id
		}

		cityInfo := CombinedCityInfo{
			ID:   finalID,
			Name: name,
		}
		cities = append(cities, cityInfo)
		cityMap[name] = cityInfo
	}

	// 2. Combine with StaticRegions
	var combinedRegions []CombinedRegionInfo
	for _, reg := range StaticRegions {
		var regCities []CombinedCityInfo
		for _, cityName := range reg.Cities {
			if cityInfo, found := cityMap[cityName]; found {
				regCities = append(regCities, cityInfo)
			}
		}
		combinedRegions = append(combinedRegions, CombinedRegionInfo{
			Name:   reg.Name,
			Cities: regCities,
		})
	}

	return c.JSON(CombinedLocationsResponse{
		Cities:  cities,
		Regions: combinedRegions,
	})
}
