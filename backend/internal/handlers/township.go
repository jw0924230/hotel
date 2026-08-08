package handlers

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"

	"hotel-backend/internal/database"
)

type TownshipCandidate struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CityName string `json:"city"`
}

type TownshipMatch struct {
	Candidate *TownshipCandidate `json:"township,omitempty"`
	Reason    string             `json:"reason,omitempty"`
}

var cityAliases = map[string][]string{
	"基隆": {"基隆市"}, "台北": {"台北市"}, "新北": {"新北市"},
	"桃園": {"桃園市"}, "新竹": {"新竹市", "新竹縣"}, "宜蘭": {"宜蘭縣"},
	"苗栗": {"苗栗縣"}, "台中": {"台中市"}, "彰化": {"彰化縣"},
	"南投": {"南投縣"}, "雲林": {"雲林縣"}, "嘉義": {"嘉義市", "嘉義縣"},
	"台南": {"台南市"}, "高雄": {"高雄市"}, "屏東": {"屏東縣"},
	"花蓮": {"花蓮縣"}, "台東": {"台東縣"}, "澎湖": {"澎湖縣"},
	"金門": {"金門縣"}, "馬祖": {"連江縣", "馬祖"},
}

func normalizeAddress(value string) string {
	value = strings.ReplaceAll(value, "臺", "台")
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '　' {
			return -1
		}
		return r
	}, value)
}

func normalizeArea(value string) string {
	value = normalizeAddress(strings.TrimSpace(value))
	for city, aliases := range cityAliases {
		if value == city {
			return city
		}
		for _, alias := range aliases {
			if value == alias {
				return city
			}
		}
	}
	return value
}

func cityFromAddress(address string) string {
	for city, aliases := range cityAliases {
		for _, alias := range aliases {
			if strings.Contains(address, alias) {
				return city
			}
		}
	}
	return ""
}

// InferTownship deterministically matches an address to a child of its city.
// It deliberately returns no match when the city is unknown or evidence conflicts.
func InferTownship(area, address string, candidates []TownshipCandidate) TownshipMatch {
	normalizedAddress := normalizeAddress(address)
	if normalizedAddress == "" {
		return TownshipMatch{Reason: "missing_address"}
	}

	areaCity := normalizeArea(area)
	addressCity := cityFromAddress(normalizedAddress)
	if areaCity == "其他" {
		areaCity = ""
	}
	if areaCity != "" && addressCity != "" && areaCity != addressCity {
		return TownshipMatch{Reason: "city_conflict"}
	}
	city := areaCity
	if city == "" {
		city = addressCity
	}
	if city == "" {
		return TownshipMatch{Reason: "city_not_found"}
	}

	matches := make([]TownshipCandidate, 0, 2)
	for _, candidate := range candidates {
		if candidate.CityName == city && strings.Contains(normalizedAddress, normalizeAddress(candidate.Name)) {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 0 {
		return TownshipMatch{Reason: "township_not_found"}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return len([]rune(matches[i].Name)) > len([]rune(matches[j].Name))
	})
	longest := len([]rune(matches[0].Name))
	if len(matches) > 1 && len([]rune(matches[1].Name)) == longest {
		return TownshipMatch{Reason: "ambiguous_township"}
	}
	return TownshipMatch{Candidate: &matches[0]}
}

type TownshipHandler struct {
	DB *database.DB
}

func NewTownshipHandler(db *database.DB) *TownshipHandler {
	return &TownshipHandler{DB: db}
}

func (h *TownshipHandler) loadCandidates(ctx context.Context) ([]TownshipCandidate, error) {
	rows, err := h.DB.Pool.Query(ctx, `
		SELECT t.id, t.name, city.name
		FROM categories t
		JOIN categories city ON city.id = t.parent_id
		WHERE t.type = 'township' AND city.type = 'city'
		ORDER BY city.sort_order, t.sort_order, t.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := make([]TownshipCandidate, 0, 368)
	for rows.Next() {
		var candidate TownshipCandidate
		if err := rows.Scan(&candidate.ID, &candidate.Name, &candidate.CityName); err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, rows.Err()
}

// Analyze suggests one township without mutating hotel data.
func (h *TownshipHandler) Analyze(c *fiber.Ctx) error {
	var input struct {
		Area    string `json:"area"`
		Address string `json:"address"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	candidates, err := h.loadCandidates(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to load townships: %v", err)})
	}
	match := InferTownship(input.Area, input.Address, candidates)
	return c.JSON(fiber.Map{
		"matched":  match.Candidate != nil,
		"township": match.Candidate,
		"reason":   match.Reason,
	})
}

type unmatchedHotel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

// Backfill analyzes only hotels whose normalized township is still NULL.
func (h *TownshipHandler) Backfill(c *fiber.Ctx) error {
	ctx := context.Background()
	candidates, err := h.loadCandidates(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to load townships: %v", err)})
	}

	rows, err := h.DB.Pool.Query(ctx, `
		SELECT id, name, COALESCE(area, ''), COALESCE(address, '')
		FROM hotels WHERE township_category_id IS NULL ORDER BY id`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to load hotels: %v", err)})
	}
	type pendingHotel struct{ id, name, area, address string }
	pending := make([]pendingHotel, 0)
	for rows.Next() {
		var hotel pendingHotel
		if err := rows.Scan(&hotel.id, &hotel.name, &hotel.area, &hotel.address); err != nil {
			rows.Close()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan hotel: %v", err)})
		}
		pending = append(pending, hotel)
	}
	rows.Close()

	tx, err := h.DB.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to begin backfill: %v", err)})
	}
	defer tx.Rollback(ctx)

	updated, skipped := 0, 0
	unmatched := make([]unmatchedHotel, 0)
	for _, hotel := range pending {
		match := InferTownship(hotel.area, hotel.address, candidates)
		if match.Candidate == nil {
			unmatched = append(unmatched, unmatchedHotel{hotel.id, hotel.name, hotel.address, match.Reason})
			continue
		}
		result, err := tx.Exec(ctx, `
			UPDATE hotels SET township_category_id = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND township_category_id IS NULL`, match.Candidate.ID, hotel.id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to update hotel %s: %v", hotel.id, err)})
		}
		if result.RowsAffected() == 1 {
			updated++
		} else {
			skipped++
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to commit backfill: %v", err)})
	}

	remaining := len(unmatched)
	details := unmatched
	truncated := false
	if len(details) > 100 {
		details = details[:100]
		truncated = true
	}
	return c.JSON(fiber.Map{
		"analyzed":          len(pending),
		"updated":           updated,
		"skipped":           skipped,
		"unmatched":         len(unmatched),
		"remaining_empty":   remaining,
		"unmatched_hotels":  details,
		"details_truncated": truncated,
	})
}
