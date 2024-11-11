package models

import (
	"time"
)
import "gorm.io/gorm"

type FilterItem struct {
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
	WatchLists       []WatchList `gorm:"foreignKey:FilterItemID"` // Optional, for reverse lookup
}

// CreateFilterItem creates a new FilterItem entry in the database
func CreateFilterItem(db *gorm.DB, filterItem FilterItem) (*FilterItem, error) {
	if err := db.Create(&filterItem).Error; err != nil {
		return nil, err
	}
	return &filterItem, nil
}

// GetFilterItemByID retrieves a FilterItem by its ID
func GetFilterItemByID(db *gorm.DB, id uint) (*FilterItem, error) {
	var filterItem FilterItem
	if err := db.First(&filterItem, id).Error; err != nil {
		return nil, err
	}
	return &filterItem, nil
}

// UpdateFilterItem updates an existing FilterItem based on its ID
func UpdateFilterItem(db *gorm.DB, id uint, updatedData FilterItem) (*FilterItem, error) {
	var filterItem FilterItem
	if err := db.First(&filterItem, id).Error; err != nil {
		return nil, err
	}

	// Update fields based on the updatedData struct
	if err := db.Model(&filterItem).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	return &filterItem, nil
}

// DeleteFilterItem deletes a FilterItem by its ID
func DeleteFilterItem(db *gorm.DB, id uint) error {
	if err := db.Delete(&FilterItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllFilterItems retrieves all FilterItem entries in the database
func GetAllFilterItems(db *gorm.DB) ([]FilterItem, error) {
	var filterItems []FilterItem
	if err := db.Find(&filterItems).Error; err != nil {
		return nil, err
	}
	return filterItems, nil
}
