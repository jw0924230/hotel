package models

import (
	"time"
)

// Category represents a general-purpose categorization record.
// The Type field distinguishes between different category kinds (e.g., "city", "hotel_category").
type Category struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
