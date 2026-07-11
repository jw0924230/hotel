package models

type HomepageHotel struct {
	City      string `json:"city"`
	SortOrder int    `json:"sort_order"`
	HotelID   string `json:"hotel_id"`
}
