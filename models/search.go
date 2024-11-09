package models

import (
	"time"
)

type SearchItem struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	PriceMin         float64    `json:"price_min"`
	PriceMax         float64    `json:"price_max"`
	City             string     `json:"city"`
	Neighborhood     string     `json:"neighborhood"`
	AreaMin          int        `json:"area_min"`
	AreaMax          int        `json:"area_max"`
	BedroomsMin      int        `json:"bedrooms_min"`
	BedroomsMax      int        `json:"bedrooms_max"`
	Category         string     `json:"category"` // rent, buy, mortgage
	AgeMin           int        `json:"age_min"`
	AgeMax           int        `json:"age_max"`
	PropertyType     string     `json:"property_type"` // apartment, villa
	FloorMin         int        `json:"floor_min"`
	FloorMax         int        `json:"floor_max"`
	HasStorage       bool       `json:"has_storage"`
	HasElevator      bool       `json:"has_elevator"`
	CreatedDateStart time.Time  `json:"created_date_start"`
	CreatedDateEnd   time.Time  `json:"created_date_end"`
	WatchLists       []WatchList `gorm:"foreignKey:SearchItemID"` // Optional, for reverse lookup
}
