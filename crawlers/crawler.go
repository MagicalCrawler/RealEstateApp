package crawlers

import (
	"context"
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
	City        City
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

// SingleCrawlerData represents metadata information for each crawler instance
type SingleCrawlerData struct {
	Successful  bool
	CPUUsage    float64
	MemoryUsage float64
	Posts       []Post
}

type Crawler interface {
	Crawl(ctx context.Context) ([]Post, error)
	CrawlPostDetails(ctx context.Context, post Post) (Post, error)
}
