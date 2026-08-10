package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"
)

// HotelHandler handles HTTP requests for hotel data.
type HotelHandler struct {
	DB *database.DB
}

// NewHotelHandler creates a new HotelHandler instance.
func NewHotelHandler(db *database.DB) *HotelHandler {
	return &HotelHandler{DB: db}
}

// List handles listing hotels with pagination and filters.
// GET /api/hotels?page=1&limit=10&area=基隆,台北&query=北極星
// The area parameter now supports comma-separated values for multi-city filtering.
func (h *HotelHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	areaRaw := c.Query("area", "")
	townshipRaw := c.Query("township_id", "")
	townshipsRaw := c.Query("township_ids", "")
	tagRaw := c.Query("tag_id", "")
	query := c.Query("query", "")
	showDisabled := c.Query("show_disabled", "false") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	townshipIDs, err := parsePositiveIDs(townshipsRaw)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid township_ids"})
	}
	if townshipRaw != "" {
		legacyIDs, legacyErr := parsePositiveIDs(townshipRaw)
		if legacyErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid township_id"})
		}
		townshipIDs = appendUniqueIDs(townshipIDs, legacyIDs...)
	}
	var tagID *int
	if tagRaw != "" {
		parsed, err := strconv.Atoi(tagRaw)
		if err != nil || parsed < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tag_id"})
		}
		tagID = &parsed
	}

	ctx := context.Background()

	// Parse comma-separated areas into a slice
	var areas []string
	if areaRaw != "" {
		for _, a := range strings.Split(areaRaw, ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				areas = append(areas, a)
			}
		}
	}

	// Intercept homepage queries to return custom selected hotels
	if len(areas) == 1 && len(townshipIDs) == 0 && tagID == nil && limit == 6 && !showDisabled && query == "" {
		cityName := areas[0]
		checkQuery := `SELECT COUNT(*) FROM homepage_hotels WHERE city = $1`
		var customCount int
		_ = h.DB.Pool.QueryRow(ctx, checkQuery, cityName).Scan(&customCount)

		if customCount > 0 {
			selectCustomQuery := `
				SELECT h.id, h.name, COALESCE(h.area, ''), h.township_category_id,
				       COALESCE(t.name, ''), COALESCE(h.address, ''),
				       COALESCE(h.phone, ''), COALESCE(h.category, ''),
				       COALESCE(hp.weekday_stay, 0), COALESCE(hp.holiday_stay, 0),
				       COALESCE(hp.weekday_rest_hours, 0), COALESCE(hp.weekday_rest, 0),
				       COALESCE(hp.holiday_rest_hours, 0), COALESCE(hp.holiday_rest, 0),
			       COALESCE((
			           SELECT jsonb_agg(hi.url ORDER BY hi.sort_order, hi.id)
			           FROM hotel_images hi WHERE hi.hotel_id = h.id
			       ), '[]'::jsonb),
			       COALESCE((
			           SELECT jsonb_agg(jsonb_build_object('id', tag.id, 'name', tag.name) ORDER BY ht.sort_order, tag.sort_order, tag.id)
			           FROM hotel_tags ht JOIN categories tag ON tag.id = ht.tag_id AND tag.type = 'hotel_tag'
			           WHERE ht.hotel_id = h.id
			       ), '[]'::jsonb),
			       h.is_disabled, h.created_at, h.updated_at
				FROM homepage_hotels hh
				JOIN hotels h ON h.id = hh.hotel_id
				LEFT JOIN hotel_prices hp ON hp.hotel_id = h.id
				LEFT JOIN categories t ON t.id = h.township_category_id AND t.type = 'township'
				WHERE hh.city = $1 AND h.is_disabled = FALSE
				ORDER BY hh.sort_order ASC`

			rows, err := h.DB.Pool.Query(ctx, selectCustomQuery, cityName)
			if err == nil {
				defer rows.Close()
				customHotels := []models.Hotel{}
				for rows.Next() {
					var hotel models.Hotel
					var imagesJSON []byte
					var tagsJSON []byte
					err := rows.Scan(
						&hotel.ID, &hotel.Name, &hotel.Area, &hotel.TownshipID, &hotel.Township,
						&hotel.Address, &hotel.Phone,
						&hotel.Category, &hotel.Pricing.WeekdayStay, &hotel.Pricing.HolidayStay,
						&hotel.Pricing.WeekdayRestHours, &hotel.Pricing.WeekdayRest,
						&hotel.Pricing.HolidayRestHours, &hotel.Pricing.HolidayRest,
						&imagesJSON, &tagsJSON,
						&hotel.IsDisabled, &hotel.CreatedAt, &hotel.UpdatedAt,
					)
					if err == nil {
						if len(imagesJSON) > 0 {
							_ = json.Unmarshal(imagesJSON, &hotel.Images)
						}
						if hotel.Images == nil {
							hotel.Images = []string{}
						}
						hydrateHotelTags(&hotel, tagsJSON)
						hotel.Price = formatPriceLabel(hotel.Pricing)
						customHotels = append(customHotels, hotel)
					}
				}
				if len(customHotels) > 0 {
					return c.JSON(fiber.Map{
						"data":  customHotels,
						"total": len(customHotels),
						"page":  1,
						"limit": limit,
					})
				}
			}
		}
	}

	// Base count query
	countQuery := `SELECT COUNT(*) FROM hotels h WHERE 1=1`
	if !showDisabled {
		countQuery += ` AND h.is_disabled = FALSE`
	}
	var countArgs []interface{}
	argCount := 1

	if len(areas) > 0 {
		countQuery += fmt.Sprintf(" AND h.area = ANY($%d)", argCount)
		countArgs = append(countArgs, areas)
		argCount++
	}
	if len(townshipIDs) > 0 {
		countQuery += fmt.Sprintf(" AND h.township_category_id = ANY($%d)", argCount)
		countArgs = append(countArgs, townshipIDs)
		argCount++
	}
	if tagID != nil {
		countQuery += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM hotel_tags ht WHERE ht.hotel_id = h.id AND ht.tag_id = $%d)", argCount)
		countArgs = append(countArgs, *tagID)
		argCount++
	}
	if query != "" {
		countQuery += fmt.Sprintf(" AND (h.name ILIKE $%d OR h.address ILIKE $%d)", argCount, argCount)
		countArgs = append(countArgs, "%"+query+"%")
		argCount++
	}

	var total int
	err = h.DB.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to count hotels: %v", err)})
	}

	// Base select query
	selectQuery := `
		SELECT h.id, h.name, COALESCE(h.area, ''), h.township_category_id,
		       COALESCE(t.name, ''), COALESCE(h.address, ''),
		       COALESCE(h.phone, ''), COALESCE(h.category, ''),
		       COALESCE(hp.weekday_stay, 0), COALESCE(hp.holiday_stay, 0),
		       COALESCE(hp.weekday_rest_hours, 0), COALESCE(hp.weekday_rest, 0),
		       COALESCE(hp.holiday_rest_hours, 0), COALESCE(hp.holiday_rest, 0),
	       COALESCE((
		           SELECT jsonb_agg(hi.url ORDER BY hi.sort_order, hi.id)
		           FROM hotel_images hi WHERE hi.hotel_id = h.id
	       ), '[]'::jsonb),
	       COALESCE((
	           SELECT jsonb_agg(jsonb_build_object('id', tag.id, 'name', tag.name) ORDER BY ht.sort_order, tag.sort_order, tag.id)
	           FROM hotel_tags ht JOIN categories tag ON tag.id = ht.tag_id AND tag.type = 'hotel_tag'
	           WHERE ht.hotel_id = h.id
	       ), '[]'::jsonb),
	       h.is_disabled, h.created_at, h.updated_at
		FROM hotels h
		LEFT JOIN hotel_prices hp ON hp.hotel_id = h.id
		LEFT JOIN categories t ON t.id = h.township_category_id AND t.type = 'township'
		WHERE 1=1`
	if !showDisabled {
		selectQuery += ` AND h.is_disabled = FALSE`
	}

	var selectArgs []interface{}
	argSelectCount := 1

	if len(areas) > 0 {
		selectQuery += fmt.Sprintf(" AND h.area = ANY($%d)", argSelectCount)
		selectArgs = append(selectArgs, areas)
		argSelectCount++
	}
	if len(townshipIDs) > 0 {
		selectQuery += fmt.Sprintf(" AND h.township_category_id = ANY($%d)", argSelectCount)
		selectArgs = append(selectArgs, townshipIDs)
		argSelectCount++
	}
	if tagID != nil {
		selectQuery += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM hotel_tags ht WHERE ht.hotel_id = h.id AND ht.tag_id = $%d)", argSelectCount)
		selectArgs = append(selectArgs, *tagID)
		argSelectCount++
	}
	if query != "" {
		selectQuery += fmt.Sprintf(" AND (h.name ILIKE $%d OR h.address ILIKE $%d)", argSelectCount, argSelectCount)
		selectArgs = append(selectArgs, "%"+query+"%")
		argSelectCount++
	}

	// Add order, limit, offset
	selectQuery += fmt.Sprintf(" ORDER BY h.id ASC LIMIT $%d OFFSET $%d", argSelectCount, argSelectCount+1)
	selectArgs = append(selectArgs, limit, offset)

	rows, err := h.DB.Pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotels: %v", err)})
	}
	defer rows.Close()

	hotels := []models.Hotel{}
	for rows.Next() {
		var hotel models.Hotel
		var imagesJSON []byte
		var tagsJSON []byte
		err := rows.Scan(
			&hotel.ID, &hotel.Name, &hotel.Area, &hotel.TownshipID, &hotel.Township,
			&hotel.Address, &hotel.Phone,
			&hotel.Category, &hotel.Pricing.WeekdayStay, &hotel.Pricing.HolidayStay,
			&hotel.Pricing.WeekdayRestHours, &hotel.Pricing.WeekdayRest,
			&hotel.Pricing.HolidayRestHours, &hotel.Pricing.HolidayRest,
			&imagesJSON, &tagsJSON,
			&hotel.IsDisabled, &hotel.CreatedAt, &hotel.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan hotel: %v", err)})
		}

		if len(imagesJSON) > 0 {
			_ = json.Unmarshal(imagesJSON, &hotel.Images)
		}
		if hotel.Images == nil {
			hotel.Images = []string{}
		}
		hydrateHotelTags(&hotel, tagsJSON)
		hotel.Price = formatPriceLabel(hotel.Pricing)

		hotels = append(hotels, hotel)
	}

	return c.JSON(fiber.Map{
		"data":  hotels,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Get handles retrieving a single hotel by ID.
// GET /api/hotels/:id
func (h *HotelHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}
	showDisabled := c.Query("show_disabled", "false") == "true"

	ctx := context.Background()
	var hotel models.Hotel
	var imagesJSON []byte
	var tagsJSON []byte

	query := `
		SELECT h.id, h.name, COALESCE(h.area, ''), h.township_category_id,
		       COALESCE(t.name, ''), COALESCE(h.address, ''),
		       COALESCE(h.phone, ''), COALESCE(h.fax, ''), COALESCE(h.website, ''),
		       COALESCE(h.email, ''), COALESCE(h.category, ''),
		       COALESCE(h.stay_info, ''), COALESCE(h.housing_rules, ''),
		       COALESCE(h.description, ''), COALESCE(h.booking_link, ''),
		       COALESCE(hp.weekday_stay, 0), COALESCE(hp.holiday_stay, 0),
		       COALESCE(hp.weekday_rest_hours, 0), COALESCE(hp.weekday_rest, 0),
		       COALESCE(hp.holiday_rest_hours, 0), COALESCE(hp.holiday_rest, 0),
	       COALESCE((
		           SELECT jsonb_agg(hi.url ORDER BY hi.sort_order, hi.id)
		           FROM hotel_images hi WHERE hi.hotel_id = h.id
	       ), '[]'::jsonb),
	       COALESCE((
	           SELECT jsonb_agg(jsonb_build_object('id', tag.id, 'name', tag.name) ORDER BY ht.sort_order, tag.sort_order, tag.id)
	           FROM hotel_tags ht JOIN categories tag ON tag.id = ht.tag_id AND tag.type = 'hotel_tag'
	           WHERE ht.hotel_id = h.id
	       ), '[]'::jsonb),
	       h.is_disabled, h.created_at, h.updated_at
		FROM hotels h
		LEFT JOIN hotel_prices hp ON hp.hotel_id = h.id
		LEFT JOIN categories t ON t.id = h.township_category_id AND t.type = 'township'
		WHERE h.id = $1`

	err := h.DB.Pool.QueryRow(ctx, query, id).Scan(
		&hotel.ID, &hotel.Name, &hotel.Area, &hotel.TownshipID, &hotel.Township,
		&hotel.Address, &hotel.Phone, &hotel.Fax,
		&hotel.Website, &hotel.Email, &hotel.Category, &hotel.StayInfo, &hotel.HousingRules,
		&hotel.Description, &hotel.BookingLink, &hotel.Pricing.WeekdayStay,
		&hotel.Pricing.HolidayStay, &hotel.Pricing.WeekdayRestHours, &hotel.Pricing.WeekdayRest,
		&hotel.Pricing.HolidayRestHours, &hotel.Pricing.HolidayRest,
		&imagesJSON, &tagsJSON, &hotel.IsDisabled, &hotel.CreatedAt, &hotel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "hotel not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotel: %v", err)})
	}

	if !showDisabled && hotel.IsDisabled {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "hotel not found"})
	}

	if len(imagesJSON) > 0 {
		_ = json.Unmarshal(imagesJSON, &hotel.Images)
	}
	if hotel.Images == nil {
		hotel.Images = []string{}
	}
	hydrateHotelTags(&hotel, tagsJSON)
	hotel.Price = formatPriceLabel(hotel.Pricing)

	return c.JSON(hotel)
}

// Upsert handles creating or updating a hotel.
// PUT /api/hotels/:id
func (h *HotelHandler) Upsert(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	var hotel models.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("failed to parse request body: %v", err)})
	}
	hotel.ID = id // Ensure ID in struct matches route param

	cleanTagIDs := make([]int, 0, len(hotel.TagIDs))
	seenTagIDs := make(map[int]struct{}, len(hotel.TagIDs))
	for _, tagID := range hotel.TagIDs {
		if tagID < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "旅館標籤包含無效 ID"})
		}
		if _, exists := seenTagIDs[tagID]; exists {
			continue
		}
		seenTagIDs[tagID] = struct{}{}
		cleanTagIDs = append(cleanTagIDs, tagID)
	}
	if len(cleanTagIDs) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "每間旅館最多只能設定五個標籤"})
	}
	hotel.TagIDs = cleanTagIDs

	// Validate image links (only allow meee.com.tw)
	for _, img := range hotel.Images {
		imgTrimmed := strings.TrimSpace(img)
		if imgTrimmed != "" && !strings.Contains(strings.ToLower(imgTrimmed), "meee.com.tw") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "只允許使用 meee 的圖床連結 (例如: https://i.meee.com.tw/xxxx.jpg)"})
		}
	}

	ctx := context.Background()
	if len(hotel.TagIDs) > 0 {
		var validTagCount int
		if err := h.DB.Pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM categories
			WHERE type = 'hotel_tag' AND parent_id IS NULL AND id = ANY($1)`, hotel.TagIDs,
		).Scan(&validTagCount); err != nil || validTagCount != len(hotel.TagIDs) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "旅館標籤不存在或類型不正確"})
		}
	}
	if hotel.TownshipID != nil {
		var valid bool
		err := h.DB.Pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM categories t
				JOIN categories city ON city.id = t.parent_id
				WHERE t.id = $1 AND t.type = 'township'
				  AND city.type = 'city' AND city.name = $2
			)`, *hotel.TownshipID, hotel.Area).Scan(&valid)
		if err != nil || !valid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "鄉鎮市區分類與地區不相符"})
		}
	}
	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to begin transaction: %v", err)})
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO hotels (
			id, name, area, township_category_id, address, phone, fax, website, email, category,
			stay_info, housing_rules, description, booking_link, is_disabled, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, CURRENT_TIMESTAMP
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			area = EXCLUDED.area,
			township_category_id = EXCLUDED.township_category_id,
			address = EXCLUDED.address,
			phone = EXCLUDED.phone,
			fax = EXCLUDED.fax,
			website = EXCLUDED.website,
			email = EXCLUDED.email,
			category = EXCLUDED.category,
			stay_info = EXCLUDED.stay_info,
			housing_rules = EXCLUDED.housing_rules,
			description = EXCLUDED.description,
			booking_link = EXCLUDED.booking_link,
			is_disabled = EXCLUDED.is_disabled,
			updated_at = CURRENT_TIMESTAMP
		RETURNING created_at, updated_at`

	err = tx.QueryRow(ctx, query,
		hotel.ID, hotel.Name, hotel.Area, hotel.TownshipID, hotel.Address, hotel.Phone, hotel.Fax,
		hotel.Website, hotel.Email, hotel.Category, hotel.StayInfo, hotel.HousingRules,
		hotel.Description, hotel.BookingLink, hotel.IsDisabled,
	).Scan(&hotel.CreatedAt, &hotel.UpdatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to upsert hotel: %v", err)})
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO hotel_prices (
			hotel_id, weekday_stay, holiday_stay, weekday_rest_hours,
			weekday_rest, holiday_rest_hours, holiday_rest, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (hotel_id) DO UPDATE SET
			weekday_stay = EXCLUDED.weekday_stay,
			holiday_stay = EXCLUDED.holiday_stay,
			weekday_rest_hours = EXCLUDED.weekday_rest_hours,
			weekday_rest = EXCLUDED.weekday_rest,
			holiday_rest_hours = EXCLUDED.holiday_rest_hours,
			holiday_rest = EXCLUDED.holiday_rest,
			updated_at = CURRENT_TIMESTAMP`,
		hotel.ID, hotel.Pricing.WeekdayStay, hotel.Pricing.HolidayStay,
		hotel.Pricing.WeekdayRestHours, hotel.Pricing.WeekdayRest,
		hotel.Pricing.HolidayRestHours, hotel.Pricing.HolidayRest,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to save hotel prices: %v", err)})
	}

	if _, err = tx.Exec(ctx, `DELETE FROM hotel_images WHERE hotel_id = $1`, hotel.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to replace hotel images: %v", err)})
	}

	cleanImages := make([]string, 0, len(hotel.Images))
	for _, imageURL := range hotel.Images {
		imageURL = strings.TrimSpace(imageURL)
		if imageURL == "" {
			continue
		}
		cleanImages = append(cleanImages, imageURL)
	}
	hotel.Images = cleanImages

	for sortOrder, imageURL := range hotel.Images {
		if _, err = tx.Exec(ctx,
			`INSERT INTO hotel_images (hotel_id, url, sort_order) VALUES ($1, $2, $3)`,
			hotel.ID, imageURL, sortOrder,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to save hotel image: %v", err)})
		}
	}

	if _, err = tx.Exec(ctx, `DELETE FROM hotel_tags WHERE hotel_id = $1`, hotel.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to replace hotel tags: %v", err)})
	}
	for sortOrder, tagID := range hotel.TagIDs {
		if _, err = tx.Exec(ctx, `INSERT INTO hotel_tags (hotel_id, tag_id, sort_order) VALUES ($1, $2, $3)`, hotel.ID, tagID, sortOrder); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to save hotel tag: %v", err)})
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to commit hotel changes: %v", err)})
	}
	if hotel.TownshipID != nil {
		_ = h.DB.Pool.QueryRow(ctx, `SELECT name FROM categories WHERE id = $1`, *hotel.TownshipID).Scan(&hotel.Township)
	} else {
		hotel.Township = ""
	}
	tags, err := h.loadHotelTags(ctx, hotel.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to reload hotel tags: %v", err)})
	}
	hotel.Tags = tags
	hotel.TagIDs = make([]int, 0, len(tags))
	for _, tag := range tags {
		hotel.TagIDs = append(hotel.TagIDs, tag.ID)
	}

	hotel.Price = formatPriceLabel(hotel.Pricing)
	return c.JSON(hotel)
}

func hydrateHotelTags(hotel *models.Hotel, tagsJSON []byte) {
	hotel.Tags = []models.HotelTag{}
	if len(tagsJSON) > 0 {
		_ = json.Unmarshal(tagsJSON, &hotel.Tags)
	}
	hotel.TagIDs = make([]int, 0, len(hotel.Tags))
	for _, tag := range hotel.Tags {
		hotel.TagIDs = append(hotel.TagIDs, tag.ID)
	}
}

func (h *HotelHandler) loadHotelTags(ctx context.Context, hotelID string) ([]models.HotelTag, error) {
	rows, err := h.DB.Pool.Query(ctx, `
		SELECT c.id, c.name
		FROM hotel_tags ht
		JOIN categories c ON c.id = ht.tag_id AND c.type = 'hotel_tag'
		WHERE ht.hotel_id = $1
		ORDER BY ht.sort_order ASC, c.sort_order ASC, c.id ASC`, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.HotelTag{}
	for rows.Next() {
		var tag models.HotelTag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func parsePositiveIDs(raw string) ([]int, error) {
	if strings.TrimSpace(raw) == "" {
		return []int{}, nil
	}
	ids := []int{}
	for _, value := range strings.Split(raw, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || id < 1 {
			return nil, fmt.Errorf("invalid ID")
		}
		ids = appendUniqueIDs(ids, id)
	}
	return ids, nil
}

func appendUniqueIDs(ids []int, values ...int) []int {
	seen := make(map[int]struct{}, len(ids)+len(values))
	for _, id := range ids {
		seen[id] = struct{}{}
	}
	for _, id := range values {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func formatPriceLabel(pricing models.HotelPrice) string {
	if pricing.WeekdayRest > 0 {
		if pricing.WeekdayRestHours > 0 {
			return fmt.Sprintf("休息 %dH $%d起", pricing.WeekdayRestHours, pricing.WeekdayRest)
		}
		return fmt.Sprintf("休息 $%d起", pricing.WeekdayRest)
	}
	if pricing.WeekdayStay > 0 {
		return fmt.Sprintf("住宿 $%d起", pricing.WeekdayStay)
	}
	return ""
}

// UploadImages handles uploading hotel images (limit 1MB each, max 10 images)
// POST /api/hotels/:id/upload
func (h *HotelHandler) UploadImages(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("failed to parse multipart form: %v", err)})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No images provided"})
	}

	// Validate max images (limit 10)
	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You can upload a maximum of 10 images"})
	}

	// Retrieve upload directory from environment or default to local frontend public folder
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		// Default relative path from backend/ directory to frontend public folder
		uploadDir = "../frontend/public/data/images"
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to create upload directory: %v", err)})
	}

	uploadedFilenames := []string{}

	for idx, file := range files {
		// Validate size (1MB limit)
		if file.Size > 1*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("File %s exceeds the 1MB limit", file.Filename)})
		}

		// Validate extension (simple check)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("File %s is not a supported image format (JPEG/PNG/WEBP only)", file.Filename)})
		}

		// Generate clean filename: {hotel_id}_{timestamp}_{index}{ext}
		filename := fmt.Sprintf("%s_%d_%d%s", id, time.Now().Unix(), idx, ext)
		targetPath := filepath.Join(uploadDir, filename)

		// Open source file
		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to open source file %s: %v", file.Filename, err)})
		}
		defer src.Close()

		// Create destination file
		dst, err := os.Create(targetPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to create destination file %s: %v", filename, err)})
		}
		defer dst.Close()

		// Copy content
		if _, err = io.Copy(dst, src); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to save file %s: %v", filename, err)})
		}

		uploadedFilenames = append(uploadedFilenames, filename)
	}

	ctx := context.Background()
	var currentImages []string
	rows, err := h.DB.Pool.Query(ctx,
		`SELECT url FROM hotel_images WHERE hotel_id = $1 ORDER BY sort_order, id`,
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query hotel images: %v", err)})
	}
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			rows.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan hotel image: %v", err)})
		}
		currentImages = append(currentImages, imageURL)
	}
	rows.Close()

	// Append newly uploaded filenames
	currentImages = append(currentImages, uploadedFilenames...)

	// Keep only the last 10 images if it exceeds
	if len(currentImages) > 10 {
		currentImages = currentImages[len(currentImages)-10:]
	}

	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to begin image transaction: %v", err)})
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM hotel_images WHERE hotel_id = $1`, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to replace hotel images: %v", err)})
	}
	for sortOrder, imageURL := range currentImages {
		if _, err = tx.Exec(ctx,
			`INSERT INTO hotel_images (hotel_id, url, sort_order) VALUES ($1, $2, $3)`,
			id, imageURL, sortOrder,
		); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("failed to save hotel image: %v", err)})
		}
	}
	if _, err = tx.Exec(ctx, `UPDATE hotels SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to update hotel: %v", err)})
	}
	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to commit images: %v", err)})
	}

	return c.JSON(fiber.Map{
		"message": "Images uploaded successfully",
		"images":  currentImages,
	})
}

// Helper function to handle saving multipart file (alternative method if needed)
func saveMultipartFile(file *multipart.FileHeader, dest string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
