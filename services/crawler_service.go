package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/crawlers/divar"
	"log"
	"sync"
	"time"
)

const (
	interval    = 30 * time.Minute // interval for repeating the crawler process
	numCrawlers = 1                // number of concurrent crawler instances
)

// StartCrawlers initializes and runs multiple crawler instances concurrently in a repeated cycle
func StartCrawlers() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Println("Starting crawlers...")
		var wg sync.WaitGroup

		for i := 0; i < numCrawlers; i++ {
			wg.Add(1)
			go func(crawlerID int) {
				defer wg.Done()
				log.Printf("Starting crawler #%d", crawlerID)
				runSingleCrawler()
				log.Printf("Crawler #%d completed", crawlerID)
			}(i + 1)
		}

		wg.Wait() // Wait for all crawlers to complete
		log.Println("All crawlers completed. Waiting for next cycle...")

		<-ticker.C // Wait for the next cycle
	}
}

// runSingleCrawler runs a single instance of the crawler
func runSingleCrawler() {
	ctx := context.Background()
	crawler := divar.NewDivarRealEstateCrawler()

	// Example: crawl top cities with a page limit
	posts, err := crawler.Crawl(ctx, 1)
	if err != nil {
		log.Println("Faced error on crawl =>", err)
	}

	for _, post := range posts {
		// Convert post to JSON
		postJSON, err := json.MarshalIndent(post, "", "  ")
		if err != nil {
			log.Fatalf("error converting post to JSON: %v", err)
		}
		fmt.Println(string(postJSON))
	}
}
