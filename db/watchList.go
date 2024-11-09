package models

import (
	"time"

	"gorm.io/gorm"
)

func CreateWatchList(db *gorm.DB, userID uint, searchItem SearchItem, interval int) (*WatchList, error) {
	// Save the SearchItem first
	if err := db.Create(&searchItem).Error; err != nil {
		return nil, err
	}

	// Create the WatchList with the associated SearchItemID
	watchList := WatchList{
		UserID:          userID,
		SearchItemID:    searchItem.ID,
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
