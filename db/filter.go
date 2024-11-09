package models

import "gorm.io/gorm"

// CreateFilterItem creates a new FilterItem entry in the database
func CreateFilterItem(db *gorm.DB, filterItem FilterItem) (*FilterItem, error) {
	if err := db.Create(&filterItem).Error; err != nil {
		return nil, err
	}
	return &filterItem, nil
}
