package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hotel-backend/internal/config"
	"hotel-backend/internal/database"
)

// HotelSummary holds summary info from area files
type HotelSummary struct {
	Area  string
	Price string
}

type AreaHotelItem struct {
	Name  string `json:"name"`
	Area  string `json:"area"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

var priceNumberPattern = regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)

func parsePriceNumbers(value string) []float64 {
	matches := priceNumberPattern.FindAllString(value, -1)
	numbers := make([]float64, 0, len(matches))
	for _, match := range matches {
		number, err := strconv.ParseFloat(match, 64)
		if err == nil {
			numbers = append(numbers, number)
		}
	}
	return numbers
}

func main() {
	log.Println("🔄 Starting enhanced hotel data import from areas and details JSON files...")

	// 1. Load configuration
	cfg := config.Load()

	// 2. Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Connected to database")

	// Ensure the schema is initialized (this will run table creation and column type alters)
	if err := db.InitSchema(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize schema: %v", err)
	}
	log.Println("✅ Schema initialized/verified")

	// 3. Detect directory paths
	// Running from backend/cmd/import_hotels/ or /backend/ or root
	frontendDir := "../frontend"
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "frontend"
	}
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		frontendDir = "../../frontend"
	}

	areasDir := filepath.Join(frontendDir, "app/data/areas")
	detailsDir := filepath.Join(frontendDir, "data_json/hotels/details")

	log.Printf("📂 Areas directory: %s", areasDir)
	log.Printf("📂 Details directory: %s", detailsDir)

	// 4. Parse all area files to build a map: id -> HotelSummary
	summaryMap := make(map[string]HotelSummary)
	areaFiles, err := ioutil.ReadDir(areasDir)
	if err == nil {
		log.Printf("⏳ Processing area summary files...")
		for _, f := range areaFiles {
			if !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			path := filepath.Join(areasDir, f.Name())
			data, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			var items []AreaHotelItem
			if err := json.Unmarshal(data, &items); err != nil {
				continue
			}

			for _, item := range items {
				// Extract ID from link (e.g. "https://www.qk.to/inquiry.asp?n=1264" -> "1264")
				idx := strings.Index(item.Link, "n=")
				if idx == -1 {
					continue
				}
				id := item.Link[idx+2:]
				if id == "" {
					continue
				}

				summaryMap[id] = HotelSummary{
					Area:  item.Area,
					Price: item.Price,
				}
			}
		}
		log.Printf("✅ Loaded %d hotel summaries from area JSONs", len(summaryMap))
	} else {
		log.Printf("⚠️ Warning: Failed to read areas directory: %v. Area names and summary prices will not be merged from area JSONs.", err)
	}

	// 5. Scan details JSON folder
	files, err := ioutil.ReadDir(detailsDir)
	if err != nil {
		log.Fatalf("❌ Failed to read details JSON directory: %v", err)
	}

	log.Printf("⏳ Importing hotel details...")

	count := 0
	skipped := 0
	imported := 0

	query := `
		INSERT INTO hotels (
			id, name, area, address, phone, fax, website, email, category, 
			stay_info, housing_rules, description, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			area = EXCLUDED.area,
			address = EXCLUDED.address,
			phone = EXCLUDED.phone,
			fax = EXCLUDED.fax,
			website = EXCLUDED.website,
			email = EXCLUDED.email,
			category = EXCLUDED.category,
			stay_info = EXCLUDED.stay_info,
			housing_rules = EXCLUDED.housing_rules,
			description = EXCLUDED.description,
			updated_at = CURRENT_TIMESTAMP`

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(detailsDir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("⚠️ Failed to read file %s: %v", file.Name(), err)
			skipped++
			continue
		}

		var hotelMap map[string]interface{}
		if err := json.Unmarshal(data, &hotelMap); err != nil {
			log.Printf("⚠️ Failed to parse JSON in file %s: %v", file.Name(), err)
			skipped++
			continue
		}

		id := toString(hotelMap["id"])
		name := toString(hotelMap["name"])

		if id == "" || name == "" {
			log.Printf("⚠️ Skipped file %s due to empty id or name", file.Name())
			skipped++
			continue
		}

		// Retrieve correct area and price summary from area files map
		area := toString(hotelMap["area"])
		price := toString(hotelMap["price"])

		if summary, exists := summaryMap[id]; exists {
			if area == "" {
				area = summary.Area
			}
			if price == "" {
				price = summary.Price
			}
		}

		address := toString(hotelMap["address"])
		phone := toString(hotelMap["phone"])
		fax := toString(hotelMap["fax"])
		website := toString(hotelMap["website"])
		email := toString(hotelMap["email"])
		category := toString(hotelMap["category"])
		stayInfo := toString(hotelMap["stay_info"])
		housingRules := toString(hotelMap["housing_rules"])
		priceAccommodation := toString(hotelMap["price_accommodation"])
		priceRest := toString(hotelMap["price_rest"])
		description := toString(hotelMap["description"])

		_, err = db.Pool.Exec(context.Background(), query,
			id, name, area, address, phone, fax, website, email, category,
			stayInfo, housingRules, description,
		)

		if err != nil {
			log.Printf("❌ Failed to insert hotel %s (%s): %v", name, id, err)
			skipped++
		} else {
			stayNumbers := parsePriceNumbers(priceAccommodation)
			restNumbers := parsePriceNumbers(priceRest)
			var weekdayStay, holidayStay, weekdayRest, holidayRest int
			var restHours float64
			if len(stayNumbers) > 0 {
				weekdayStay = int(stayNumbers[0])
			}
			if len(stayNumbers) > 1 {
				holidayStay = int(stayNumbers[1])
			}
			if len(restNumbers) >= 3 {
				restHours = restNumbers[0]
				weekdayRest = int(restNumbers[1])
				if strings.Contains(priceRest, "假日") {
					holidayRest = int(restNumbers[2])
				}
			} else {
				if len(restNumbers) > 0 {
					weekdayRest = int(restNumbers[0])
				}
				if len(restNumbers) > 1 && strings.Contains(priceRest, "假日") {
					holidayRest = int(restNumbers[1])
				}
			}

			_, priceErr := db.Pool.Exec(context.Background(), `
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
				id, weekdayStay, holidayStay, int(restHours), weekdayRest, int(restHours), holidayRest,
			)
			if priceErr != nil {
				log.Printf("⚠️ Failed to import prices for hotel %s (%s): %v", name, id, priceErr)
			}
			imported++
		}

		count++
		if count%500 == 0 {
			log.Printf("⏳ Processed %d/%d files...", count, len(files))
		}
	}

	log.Printf("🎉 Enhanced Import Complete! Total processed: %d, Imported/Upserted: %d, Skipped/Failed: %d", count, imported, skipped)
}
