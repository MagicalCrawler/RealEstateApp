package models

import "gorm.io/gorm"

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
