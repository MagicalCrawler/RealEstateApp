package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/crawlers"
	"github.com/MagicalCrawler/RealEstateApp/crawlers/divar"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CrawlerService manages the crawling process
type CrawlerService struct {
	crawlers    []crawlers.Crawler
	cityService *CityService
}

// NewCrawlerService creates a new instance of CrawlerService
func NewCrawlerService() *CrawlerService {
	return &CrawlerService{
		crawlers: []crawlers.Crawler{
			divar.NewDivarCrawler(),
			// Add other crawler implementations here
		},
		cityService: NewCityService(),
	}
}

// Start begins the crawling process
func (s *CrawlerService) Start() {
	go s.run()
}

// run executes the crawling cycle at regular intervals
func (s *CrawlerService) run() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		s.executeCrawlCycle()
		<-ticker.C
	}
}

// executeCrawlCycle performs a single crawling cycle
func (s *CrawlerService) executeCrawlCycle() {
	cities, err := s.cityService.GetCities()
	if err != nil {
		log.Printf("Failed to get cities: %v", err)
		return
	}

	chunkedCities := chunkCities(cities, 10)
	var posts []crawlerModels.Post

	var (
		totalCPU    float64
		totalMemory float64
	)

	startTime := time.Now()

	for _, cityChunk := range chunkedCities {
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())

		var avgCPU, avgMemory float64
		monitorDone := make(chan struct{})
		go func() {
			defer close(monitorDone)
			avgCPU, avgMemory, _ = monitorResources(ctx, 2*time.Second)
		}()

		for _, crawler := range s.crawlers {
			for _, city := range cityChunk {
				wg.Add(1)
				go func(crawler crawlers.Crawler, city crawlerModels.City) {
					defer wg.Done()
					result, err := crawler.Crawl(ctx, city)
					if err != nil {
						log.Printf("Crawler error for city %s: %v", city.Name, err)
						return
					}
					posts = append(posts, result...)
				}(crawler, city)
			}
		}

		wg.Wait()
		cancel()
		<-monitorDone

		totalCPU += avgCPU
		totalMemory += avgMemory

		log.Println("Chunk completed. Moving to next chunk...")
		time.Sleep(5 * time.Second)
	}

	endTime := time.Now()
	executionTime := endTime.Sub(startTime)

	session := CrawlerSession{
		StartTime:     startTime,
		EndTime:       endTime,
		ExecutionTime: executionTime,
		TotalCPU:      totalCPU / float64(len(chunkedCities)),
		TotalMemory:   totalMemory / float64(len(chunkedCities)),
		Posts:         posts,
	}

	// Save session data if needed
	session.saveToJSONFile("file.json")
	log.Println("All crawlers completed. Waiting for next cycle...")
}

// Helper functions and types

// CrawlerSession represents a crawling session with resource usage stats
type CrawlerSession struct {
	StartTime     time.Time
	EndTime       time.Time
	ExecutionTime time.Duration
	TotalCPU      float64
	TotalMemory   float64
	Posts         []crawlerModels.Post
}

// chunkCities splits the cities into smaller chunks
func chunkCities(cities []crawlerModels.City, chunkSize int) [][]crawlerModels.City {
	var chunks [][]crawlerModels.City
	for i := 0; i < len(cities); i += chunkSize {
		end := i + chunkSize
		if end > len(cities) {
			end = len(cities)
		}
		chunks = append(chunks, cities[i:end])
	}
	return chunks
}

// monitorResources monitors CPU and memory usage
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

// calculateAverage calculates the average of a slice of float64 numbers
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

// SaveToJSONFile saves the CrawlerSession to a JSON file
func (cs *CrawlerSession) saveToJSONFile(filename string) error {
	jsonData, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}
