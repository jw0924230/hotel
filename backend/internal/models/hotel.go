package models

import (
	"time"
)

type HotelPrice struct {
	WeekdayStay      int `json:"weekday_stay"`
	HolidayStay      int `json:"holiday_stay"`
	WeekdayRestHours int `json:"weekday_rest_hours"`
	WeekdayRest      int `json:"weekday_rest"`
	HolidayRestHours int `json:"holiday_rest_hours"`
	HolidayRest      int `json:"holiday_rest"`
}

// Hotel represents the hotel database structure and JSON mapping
type Hotel struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Area         string     `json:"area"`
	TownshipID   *int       `json:"township_id"`
	Township     string     `json:"township"`
	Address      string     `json:"address"`
	Phone        string     `json:"phone"`
	Fax          string     `json:"fax"`
	Website      string     `json:"website"`
	Email        string     `json:"email"`
	Category     string     `json:"category"`
	StayInfo     string     `json:"stay_info"`
	HousingRules string     `json:"housing_rules"`
	Pricing      HotelPrice `json:"pricing"`
	Price        string     `json:"price"` // Derived compatibility label for front-end lists.
	Description  string     `json:"description"`
	BookingLink  string     `json:"booking_link"`
	IsDisabled   bool       `json:"is_disabled"`
	Images       []string   `json:"images"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
