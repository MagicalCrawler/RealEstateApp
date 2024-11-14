package crawlers

import (
	"context"
)

type Post struct {
	ID                  string
	Title               string
	Price               string
	Link                string
	Images              []string
	Description         string
	Area                string
	YearBuilt           string
	Rooms               string
	PricePerSquareMeter string
	TotalPrice          string
	Floor               string
	Features            []string
}

type Crawler interface {
	Crawl(ctx context.Context) ([]Post, error)
	CrawlPostDetails(ctx context.Context, post Post) (Post, error)
}
