package models

import "time"

type CrawlHistory struct {
	ID          uint `gorm:"primary_key;auto_increment"`
	PostNum     uint
	CpuUsage    float32 `gorm:"type:decimal(2,2)"`
	MemoryUsage uint
	RequestsNum uint
	StartedAt   time.Time
	FinishedAt  time.Time
}
