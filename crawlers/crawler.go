package crawlers

import (
	"context"
)

type Post struct {
	Title  string
	Price  string
	Link   string
	Images []string
}

type Crawler interface {
	Crawl(ctx context.Context) ([]Post, error)
	CrawlPostDetails(ctx context.Context, post Post) (Post, error)
}
