package models

import (
	"gorm.io/gorm"
	"time"
)

type WatchList struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `json:"user_id"`
	FilterItemID   uint       `json:"filter_item_id"` // Foreign key to FilterItem
	SearchItem     SearchItem `gorm:"foreignKey:FilterItemID"` // Establishes the relationship
	RefreshInterval int      `json:"refresh_interval"` // in minutes or other units
	LastChecked    time.Time `json:"last_checked"`
}