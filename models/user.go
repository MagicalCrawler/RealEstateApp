package models

import (
	"gorm.io/gorm"
)

type Role int

const (
	SUPER_ADMIN Role = iota
	ADMIN
	USER
)

type User struct {
	gorm.Model
	ID         uint   `gorm:"autoIncrement"`
	TelegramID uint64 `gorm:"uniqueIndex"`
	Role       Role
}
