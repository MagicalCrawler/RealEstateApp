package models

import (
	"gorm.io/gorm"
	"time"
)

type WatchList struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `json:"user_id"`
	SearchItemID   uint       `json:"search_item_id"` // Foreign key to SearchItem
	SearchItem     SearchItem `gorm:"foreignKey:SearchItemID"` // Establishes the relationship
	RefreshInterval int      `json:"refresh_interval"` // in minutes or other units
	LastChecked    time.Time `json:"last_checked"`
}