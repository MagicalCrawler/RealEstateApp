package models

import (
	"time"

	"gorm.io/gorm"
)

type FilterItem struct {
	ID               uint        `gorm:"primaryKey" json:"id"`
	PriceMin         float64     `json:"price_min"`
	PriceMax         float64     `json:"price_max"`
	City             string      `json:"city"`
	Neighborhood     string      `json:"neighborhood"`
	AreaMin          int         `json:"area_min"`
	AreaMax          int         `json:"area_max"`
	BedroomsMin      int         `json:"bedrooms_min"`
	BedroomsMax      int         `json:"bedrooms_max"`
	Category         string      `json:"category"` // rent, buy, mortgage
	AgeMin           int         `json:"age_min"`
	AgeMax           int         `json:"age_max"`
	PropertyType     string      `json:"property_type"` // apartment, villa
	FloorMin         int         `json:"floor_min"`
	FloorMax         int         `json:"floor_max"`
	HasStorage       bool        `json:"has_storage"`
	HasElevator      bool        `json:"has_elevator"`
	CreatedDateStart time.Time   `json:"created_date_start"`
	CreatedDateEnd   time.Time   `json:"created_date_end"`
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

func SearchPostHistory(db *gorm.DB, filter FilterItem) ([]PostHistory, error) {
	var posts []PostHistory

	query := db.Model(&PostHistory{})

	// Apply filters based on FilterItem fields
	if filter.PriceMin > 0 {
		query = query.Where("price >= ?", filter.PriceMin)
	}
	if filter.PriceMax > 0 {
		query = query.Where("price <= ?", filter.PriceMax)
	}
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.Neighborhood != "" {
		query = query.Where("neighbor = ?", filter.Neighborhood)
	}
	if filter.AreaMin > 0 {
		query = query.Where("area >= ?", filter.AreaMin)
	}
	if filter.AreaMax > 0 {
		query = query.Where("area <= ?", filter.AreaMax)
	}
	if filter.BedroomsMin > 0 {
		query = query.Where("bedroom_num >= ?", filter.BedroomsMin)
	}
	if filter.BedroomsMax > 0 {
		query = query.Where("bedroom_num <= ?", filter.BedroomsMax)
	}
	if filter.Category != "" {
		query = query.Where("buy_mode = ?", filter.Category)
	}
	if filter.AgeMin > 0 {
		query = query.Where("age >= ?", filter.AgeMin)
	}
	if filter.AgeMax > 0 {
		query = query.Where("age <= ?", filter.AgeMax)
	}
	if filter.PropertyType != "" {
		query = query.Where("building = ?", filter.PropertyType)
	}
	if filter.FloorMin > 0 {
		query = query.Where("floors_num >= ?", filter.FloorMin)
	}
	if filter.FloorMax > 0 {
		query = query.Where("floors_num <= ?", filter.FloorMax)
	}
	if filter.HasStorage {
		query = query.Where("has_ware_house = ?", filter.HasStorage)
	}
	if filter.HasElevator {
		query = query.Where("has_elevator = ?", filter.HasElevator)
	}
	if !filter.CreatedDateStart.IsZero() {
		query = query.Where("created_at >= ?", filter.CreatedDateStart)
	}
	if !filter.CreatedDateEnd.IsZero() {
		query = query.Where("created_at <= ?", filter.CreatedDateEnd)
	}

	// Execute query
	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
