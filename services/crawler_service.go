package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/crawlers"
	"github.com/MagicalCrawler/RealEstateApp/crawlers/divar"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"os"
	"sync"
	"time"
)

const (
	interval       = 30 * time.Minute
	sampleInterval = 2 * time.Second
	maxRetries     = 3
)

func ChunkCities(cities []crawlers.City, chunkSize int) [][]crawlers.City {
	var chunks [][]crawlers.City
	for i := 0; i < len(cities); i += chunkSize {
		end := i + chunkSize
		if end > len(cities) {
			end = len(cities)
		}
		chunks = append(chunks, cities[i:end])
	}
	return chunks
}

func StartCrawlers() {
	cities, err := fetchCitiesFromAPIWithCache()
	if err != nil {
		log.Fatalf("Could not fetch cities: %v", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Println("Starting crawlers...")

		chunkedCities := ChunkCities(cities, 10)

		for _, cityChunk := range chunkedCities {
			var wg sync.WaitGroup

			for i, city := range cityChunk {
				wg.Add(1)
				go func(crawlerID int, city crawlers.City) {
					defer wg.Done()
					log.Printf("Starting crawler #%d for city %s", crawlerID, city.Name)
					crawlerMetadata := runSingleCrawler(crawlerID, city)
					saveMetadataToFile(crawlerMetadata)
					log.Printf("Crawler #%d completed for city %s", crawlerID, city.Name)
				}(i+1, city)
			}

			wg.Wait()
			log.Println("Chunk completed. Moving to next chunk...")

			time.Sleep(5 * time.Second)
		}
		log.Println("All crawlers completed. Waiting for next cycle...")
		<-ticker.C
	}
}

func monitorResources(ctx context.Context, sampleInterval time.Duration) (float64, float64, error) {
	var cpuSamples, memSamples []float64
	ticker := time.NewTicker(sampleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return calculateAverage(cpuSamples), calculateAverage(memSamples), nil

		case <-ticker.C:
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				log.Printf("Error getting CPU usage: %v", err)
				continue
			}
			cpuSamples = append(cpuSamples, cpuPercent[0])

			memStat, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Error getting memory usage: %v", err)
				continue
			}
			memSamples = append(memSamples, memStat.UsedPercent)
		}
	}
}

func calculateAverage(samples []float64) float64 {
	var sum float64
	for _, sample := range samples {
		sum += sample
	}
	if len(samples) == 0 {
		return 0
	}
	return sum / float64(len(samples))
}

func runSingleCrawler(crawlerID int, city crawlers.City) crawlers.CrawlerMetadata {
	var (
		posts       []crawlers.Post
		successful  bool
		attempts    int
		totalCPU    float64
		totalMemory float64
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cityURL := fmt.Sprintf("https://divar.ir/s/%s/real-estate", city.Slug)
	startTime := time.Now()

	for attempts = 0; attempts < maxRetries; attempts++ {
		resourceCtx, resourceCancel := context.WithCancel(ctx)

		var (
			averageCPUUsage, averageMemoryUsage float64
			monitorErr                          error
		)

		go func() {
			averageCPUUsage, averageMemoryUsage, monitorErr = monitorResources(resourceCtx, sampleInterval)
		}()

		crawler := divar.NewDivarRealEstateCrawler(cityURL)
		result, err := crawler.Crawl(ctx, 1)
		resourceCancel()

		if monitorErr != nil {
			log.Printf("Error during resource monitoring for attempt %d: %v", attempts+1, monitorErr)
		}

		if err != nil {
			log.Printf("Error crawling city %s on attempt %d: %v", city.Name, attempts+1, err)
			continue
		}

		if len(result) > 0 {
			successful = true
			posts = append(posts, result...)
			totalCPU += averageCPUUsage
			totalMemory += averageMemoryUsage
			break
		}

		log.Printf("No posts found for city %s on attempt %d, retrying...", city.Name, attempts+1)
		totalCPU += averageCPUUsage
		totalMemory += averageMemoryUsage
		time.Sleep(5 * time.Second)
	}

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	averageCPUUsage := totalCPU / float64(attempts)
	averageMemoryUsage := totalMemory / float64(attempts)

	return crawlers.CrawlerMetadata{
		CrawlerID:     crawlerID,
		City:          city,
		StartTime:     startTime,
		EndTime:       endTime,
		ExecutionTime: executionTime,
		CPUUsage:      averageCPUUsage,
		MemoryUsage:   averageMemoryUsage,
		Posts:         posts,
		Successful:    successful,
	}
}

func saveMetadataToFile(metadata crawlers.CrawlerMetadata) {
	filename := fmt.Sprintf("crawler_metadata_%s_%d.json", metadata.City.Slug, metadata.CrawlerID)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create metadata file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		log.Fatalf("Failed to write metadata to file: %v", err)
	}
}
