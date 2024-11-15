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

type CrawlerSession struct {
	TotalPosts      []crawlers.Post
	TotalCPU        float64
	TotalMemory     float64
	StartTime       time.Time
	EndTime         time.Time
	ExecutionTime   time.Duration
	SuccessfulCount int
}

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
		session := CrawlerSession{}

		chunkedCities := ChunkCities(cities, 10)

		startTime := time.Now()

		for _, cityChunk := range chunkedCities {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())

			var avgCPU, avgMemory float64
			monitorDone := make(chan struct{})
			go func() {
				defer close(monitorDone)
				avgCPU, avgMemory, _ = monitorResources(ctx, sampleInterval)
			}()

			for _, city := range cityChunk {
				wg.Add(1)
				go func(city crawlers.City) {
					defer wg.Done()
					log.Printf("Starting crawler for city %s", city.Name)
					crawlerMetadata := runSingleCrawler(city)
					if crawlerMetadata.Successful {
						session.TotalPosts = append(session.TotalPosts, crawlerMetadata.Posts...)
						session.SuccessfulCount++
					}
					log.Printf("Crawler; completed for city %s", city.Name)
				}(city)
			}

			wg.Wait()

			cancel()
			<-monitorDone

			session.TotalCPU += avgCPU
			session.TotalMemory += avgMemory

			log.Println("Chunk completed. Moving to next chunk...")
			time.Sleep(5 * time.Second)
		}

		endTime := time.Now()
		executionTime := endTime.Sub(startTime)

		session.StartTime = startTime
		session.EndTime = endTime
		session.ExecutionTime = executionTime

		session.TotalCPU /= float64(len(chunkedCities))
		session.TotalMemory /= float64(len(chunkedCities))

		saveSessionData(session)
		log.Println("All crawlers completed. Waiting for next cycle...")
		<-ticker.C
	}
}

func monitorResources(ctx context.Context, sampleInterval time.Duration) (float64, float64, error) {
	var (
		cpuSamples []float64
		memSamples []float64
	)

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
			if len(cpuPercent) > 0 {
				cpuSamples = append(cpuSamples, cpuPercent[0])
			}

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

func runSingleCrawler(city crawlers.City) crawlers.SingleCrawlerData {
	var (
		posts      []crawlers.Post
		successful bool
		attempts   int
	)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cityURL := fmt.Sprintf("https://divar.ir/s/%s/real-estate", city.Slug)

	for attempts = 0; attempts < maxRetries; attempts++ {
		crawler := divar.NewDivarRealEstateCrawler(cityURL)
		result, err := crawler.Crawl(2, city)

		if err != nil {
			log.Printf("Error crawling city %s on attempt %d: %v", city.Name, attempts+1, err)
			continue
		}

		if len(result) > 0 {
			successful = true
			posts = append(posts, result...)
			break
		}

		log.Printf("No posts found for city %s on attempt %d, retrying...", city.Name, attempts+1)
		time.Sleep(5 * time.Second)
	}

	return crawlers.SingleCrawlerData{
		Posts:      posts,
		Successful: successful,
	}
}

func saveSessionData(session CrawlerSession) {
	filename := fmt.Sprintf("crawler_session_%d.json", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create session file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(session); err != nil {
		log.Fatalf("Failed to write session data to file: %v", err)
	}
}
