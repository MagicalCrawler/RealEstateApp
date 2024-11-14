package crawlers

import (
	"context"
	"time"
)

type RentalMetadata struct {
	Capacity        string
	NormalDayPrice  string
	WeekendPrice    string
	HolidayPrice    string
	ExtraPersonCost string
}

type Post struct {
	ID          string
	Title       string
	Price       string
	Link        string
	Images      []string
	Description string
	Area        string
	YearBuilt   string
	Rooms       string

	PricePerSquareMeter string
	TotalPrice          string
	Floor               string
	Features            []string

	Deposit     string
	MonthlyRent string

	DepositOnRentDesc string
	RentalMetadata    *RentalMetadata
}

type City struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Level string `json:"level"`
}

type CityResponse struct {
	Cities []City `json:"cities"`
}

// CrawlerMetadata represents metadata information for each crawler instance
type CrawlerMetadata struct {
	CrawlerID     int
	City          City
	Successful    bool
	StartTime     time.Time
	EndTime       time.Time
	ExecutionTime time.Duration
	CPUUsage      float64
	MemoryUsage   float64
	Posts         []Post
}

type Crawler interface {
	Crawl(ctx context.Context) ([]Post, error)
	CrawlPostDetails(ctx context.Context, post Post) (Post, error)
}
