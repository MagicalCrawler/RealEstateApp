package models

import (
	"time"
)

type WatchList struct {
	ID              uint `gorm:"primaryKey" json:"id"`
	UserID          uint `json:"user_id"`
	User            User
	FilterItemID    uint         `json:"filter_item_id"`          // Foreign key to FilterItem
	FilterItem      FilterItem   `gorm:"foreignKey:FilterItemID"` // Establishes the relationship
	RefreshInterval int          `json:"refresh_interval"`        // in minutes or other units
	LastChecked     time.Time    `json:"last_checked"`
	FilterItems     []FilterItem `gorm:"foreignKey:UserID"` // Reverse relationship: a user has many filter items
}
