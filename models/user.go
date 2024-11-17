package models

import (
	"gorm.io/gorm"
)

type Role int
type UserType int

const (
	SUPER_ADMIN Role = iota
	ADMIN
	USER
)
const (
	FREE UserType = iota
	PREMIUM
)

type User struct {
	gorm.Model
	ID         uint   `gorm:"autoIncrement"`
	TelegramID uint64 `gorm:"uniqueIndex"`
	Role       Role
	Type       UserType
	FilterItems     []FilterItem `gorm:"foreignKey:UserID"` // Reverse relationship: a user has many filter items
}
