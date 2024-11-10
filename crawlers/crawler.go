package crawlers

import (
	"context"
)

type Post struct {
	Title string
	Price string
	Link  string
	Image string
}

type Crawler interface {
	Crawl(ctx context.Context) ([]Post, error)
	CrawlPostDetails(ctx context.Context, post Post) (Post, error)
}
