package models

import "time"

type Post struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Title     string `gorm:"type:text"`
	DetailID  uint   `gorm:"foreignKey:HistoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
