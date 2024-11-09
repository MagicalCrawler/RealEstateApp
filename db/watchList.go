package models

func CreateWatchList(db *gorm.DB, userID uint, searchItem SearchItem, interval int) (*WatchList, error) {
	// Save the SearchItem first
	if err := db.Create(&searchItem).Error; err != nil {
		return nil, err
	}

	// Create the WatchList with the associated SearchItemID
	watchList := WatchList{
		UserID:         userID,
		SearchItemID:   searchItem.ID,
		RefreshInterval: interval,
		LastChecked:    time.Now(),
	}

	if err := db.Create(&watchList).Error; err != nil {
		return nil, err
	}
	return &watchList, nil
}
