package crawlers

import (
	"context"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
)

// Crawler is the interface that all crawler implementations must satisfy
type Crawler interface {
	Crawl(ctx context.Context, city crawlerModels.City) ([]crawlerModels.Post, error)
	CrawlPostDetails(ctx context.Context, postURL string) (crawlerModels.Post, error)
}
