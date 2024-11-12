package models

import (
	"gorm.io/gorm"
	"time"
)

type WatchList struct {
	ID              uint `gorm:"primaryKey" json:"id"`
	UserID          uint `json:"user_id"`
	User            User
	FilterItemID    uint       `json:"filter_item_id"`          // Foreign key to FilterItem
	FilterItem      FilterItem `gorm:"foreignKey:FilterItemID"` // Establishes the relationship
	RefreshInterval int        `json:"refresh_interval"`        // in minutes or other units
	LastChecked     time.Time  `json:"last_checked"`
}

func CreateWatchList(db *gorm.DB, userID uint, filterItem FilterItem, interval int) (*WatchList, error) {
	// Save the FilterItem first
	if err := db.Create(&filterItem).Error; err != nil {
		return nil, err
	}

	// Create the WatchList with the associated FilterItemID
	watchList := WatchList{
		UserID:          userID,
		FilterItemID:    filterItem.ID,
		RefreshInterval: interval,
		LastChecked:     time.Now(),
	}

	if err := db.Create(&watchList).Error; err != nil {
		return nil, err
	}
	return &watchList, nil
}

func GetAllWatchListsForUser(db *gorm.DB, userID uint) ([]WatchList, error) {
	var watchLists []WatchList
	if err := db.Preload("FilterItem").Where("user_id = ?", userID).Find(&watchLists).Error; err != nil {
		return nil, err
	}
	return watchLists, nil
}

func DeleteWatchList(db *gorm.DB, watchListID uint) error {
	var watchList WatchList
	if err := db.First(&watchList, watchListID).Error; err != nil {
		return err
	}

	// Delete the associated FilterItem
	if err := db.Delete(&FilterItem{}, watchList.FilterItemID).Error; err != nil {
		return err
	}

	// Delete the WatchList entry
	if err := db.Delete(&WatchList{}, watchListID).Error; err != nil {
		return err
	}

	return nil
}

func UpdateWatchList(db *gorm.DB, watchListID uint, updatedData WatchList, updatedFilterItem FilterItem) (*WatchList, error) {
	var watchList WatchList
	if err := db.First(&watchList, watchListID).Error; err != nil {
		return nil, err
	}

	// Update FilterItem fields if needed
	updatedFilterItem.ID = watchList.FilterItemID // Ensure we're updating the correct FilterItem
	if err := db.Model(&FilterItem{}).Where("id = ?", watchList.FilterItemID).Updates(updatedFilterItem).Error; err != nil {
		return nil, err
	}

	// Update WatchList fields
	if err := db.Model(&watchList).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	return &watchList, nil
}

func GetWatchListByID(db *gorm.DB, watchListID uint) (*WatchList, error) {
	var watchList WatchList
	if err := db.Preload("FilterItem").First(&watchList, watchListID).Error; err != nil {
		return nil, err
	}
	return &watchList, nil
}

// GetAllWatchLists retrieves all watch list entries along with their associated filter criteria
func GetAllWatchLists(db *gorm.DB) ([]WatchList, error) {
	var watchLists []WatchList
	// Use Preload to load associated FilterItem data for each WatchList entry
	if err := db.Preload("FilterItem").Find(&watchLists).Error; err != nil {
		return nil, err
	}
	return watchLists, nil
}
