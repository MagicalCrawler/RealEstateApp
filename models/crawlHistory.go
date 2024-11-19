package models

import (
	"gorm.io/gorm"
	"time"
)

type CrawlHistory struct {
	PostNum     uint
	CpuUsage    float32 `gorm:"type:decimal(7,2)"`
	MemoryUsage float32 `gorm:"type:decimal(7,2)"`
	RequestsNum uint
	StartedAt   time.Time
	FinishedAt  time.Time
	gorm.Model
}
