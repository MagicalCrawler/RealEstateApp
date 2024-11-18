package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/crawlers"
	"github.com/MagicalCrawler/RealEstateApp/crawlers/divar"
	"github.com/MagicalCrawler/RealEstateApp/crawlers/sheypoor"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"io/ioutil"
	"log"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CrawlerService manages the crawling process
type CrawlerService struct {
	crawlers    []crawlers.Crawler
	cityService *CityService
	repository  *db.PostRepo
	logger      *slog.Logger
}

// NewCrawlerService creates a new instance of CrawlerService
func NewCrawlerService(repository *db.PostRepo) *CrawlerService {
	return &CrawlerService{
		crawlers: []crawlers.Crawler{
			divar.NewDivarCrawler(),
			sheypoor.NewSheypoorCrawler(),
		},
		cityService: NewCityService(),
		repository:  repository,
		logger:      utils.NewLogger("CrawlerService"),
	}
}

// Start begins the crawling process
func (s *CrawlerService) Start() {
	go s.run()
}

// run executes the crawling cycle at regular intervals
func (s *CrawlerService) run() {
	jobTimer, err := strconv.Atoi(utils.GetConfig("CRAWLER_INTERVAL"))
	if err != nil {
		jobTimer = 30
	}

	ticker := time.NewTicker(time.Duration(jobTimer) * time.Minute)
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
		s.logger.Error("Failed to get cities", err)
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
						s.logger.Error("Failed to crawl city", city.Name, ", Error: ", err)
						return
					}

					// تبدیل اعداد فارسی به انگلیسی در داده‌های بازگشتی
					for i := range result {
						result[i] = processPost(result[i])
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

		slog.Info("Chunk completed. Moving to next chunk...")
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
	err = session.saveToJSONFile("file.json")
	if err != nil {
		fmt.Println("Error saving session to JSON:", err)
	}

	err = mapAndSaveCrawlerSession(session, *s.repository)
	if err != nil {
		s.logger.Error("Error saving crawler session:", err)
	}

	s.logger.Info("All crawlers completed. Waiting for next cycle...")
}

// Helper functions and types

// Map for replacing Arabic digits to Persian digits
// Map for replacing Persian digits to English digits
var digitMap = map[rune]rune{
	'۰': '0', '۱': '1', '۲': '2', '۳': '3',
	'۴': '4', '۵': '5', '۶': '6', '۷': '7',
	'۸': '8', '۹': '9', '٬': ',',
}

// Function to replace Persian digits with English digits
func replaceDigits(input string) string {
	var result strings.Builder
	for _, ch := range input {
		if newCh, exists := digitMap[ch]; exists {
			result.WriteRune(newCh)
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func processPost(post crawlerModels.Post) crawlerModels.Post {
	post.Title = replaceDigits(post.Title)
	post.Description = replaceDigits(post.Description)
	post.Price = replaceDigits(post.Price)
	post.TotalPrice = replaceDigits(post.TotalPrice)
	post.Deposit = replaceDigits(post.Deposit)
	post.MonthlyRent = replaceDigits(post.MonthlyRent)
	post.PricePerSquareMeter = replaceDigits(post.PricePerSquareMeter)

	// پردازش دیگر فیلدهای متنی اگر وجود دارند
	if post.RentalMetadata != nil {
		post.RentalMetadata.NormalDayPrice = replaceDigits(post.RentalMetadata.NormalDayPrice)
		post.RentalMetadata.WeekendPrice = replaceDigits(post.RentalMetadata.WeekendPrice)
		post.RentalMetadata.HolidayPrice = replaceDigits(post.RentalMetadata.HolidayPrice)
		post.RentalMetadata.ExtraPersonCost = replaceDigits(post.RentalMetadata.ExtraPersonCost)
	}

	return post
}

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

	logger := utils.NewLogger("CrawlerService")

	ticker := time.NewTicker(sampleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return calculateAverage(cpuSamples), calculateAverage(memSamples), nil
		case <-ticker.C:
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				logger.Error("Failed to get cpu percent", err)
				continue
			}
			if len(cpuPercent) > 0 {
				cpuSamples = append(cpuSamples, cpuPercent[0])
			}

			memStat, err := mem.VirtualMemory()
			if err != nil {
				logger.Error("Failed to get mem usage", err)
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

func mapAndSaveCrawlerSession(session CrawlerSession, repository db.PostRepo) error {
	// 1. نگاشت CrawlHistory
	crawlHistory := models.CrawlHistory{
		PostNum:     uint(len(session.Posts)),
		CpuUsage:    float32(math.Round(session.TotalCPU*100) / 100),
		MemoryUsage: float32(math.Round(session.TotalMemory*100) / 100),
		RequestsNum: 1, // تعداد درخواست‌ها (اگر اطلاعاتی دارید اینجا پر کنید)
		StartedAt:   session.StartTime,
		FinishedAt:  session.EndTime,
	}

	logger := utils.NewLogger("CrawlerService")

	insertedCrawlHistory, err := repository.CrawlHistorySaving(crawlHistory)
	if err != nil {
		logger.Error("failed to save CrawlHistory: ", err)
		return fmt.Errorf("failed to save CrawlHistory: %w", err)
	}

	// 2. نگاشت Posts و PostHistory
	for _, post := range session.Posts {
		// ذخیره Post
		dbPost := models.Post{
			UniqueCode: post.ID,
			Website:    "Divar", // فرض بر اینکه این اطلاعات موجود است
		}

		insertedPost, err := repository.PostSaving(dbPost.UniqueCode, "divar.ir")
		if err != nil {
			logger.Error("failed to save post: ", post.ID, "; error: ", err)
			log.Printf("failed to save post %s: %v", post.ID, err)
			continue
		}

		postHistory := models.PostHistory{
			PostID:         insertedPost.ID,
			Title:          post.Title,
			PostURL:        post.Link,
			Price:          parsePrice(post.TotalPrice),
			Deposit:        parsePrice(post.Deposit),
			Rent:           parsePrice(post.MonthlyRent),
			City:           post.City.Name,
			Neighborhood:   "",
			Area:           parseArea(post.Area),
			BedroomNum:     parseBedrooms(post.Rooms),
			Age:            parseAge(post.YearBuilt),
			FloorsNum:      parseFloors(post.Floor),
			HasStorage:     containsFeature(post.Features, "انباری"),
			HasParking:     containsFeature(post.Features, "پارکینگ"),
			HasElevator:    containsFeature(post.Features, "آسانسور"),
			ImageURL:       strings.Join(post.Images, ","),
			Description:    post.Description,
			CrawlHistoryID: insertedCrawlHistory.ID,
		}

		// بررسی وجود RentalMetadata
		if post.RentalMetadata != nil {
			postHistory.Capacity = post.RentalMetadata.Capacity
			postHistory.NormalDays = post.RentalMetadata.NormalDayPrice
			postHistory.Weekend = post.RentalMetadata.WeekendPrice
			postHistory.Holidays = post.RentalMetadata.HolidayPrice
			postHistory.CostPerPerson = post.RentalMetadata.ExtraPersonCost
		}

		_, err = repository.PostHistorySaving(postHistory, insertedPost, insertedCrawlHistory)
		if err != nil {
			logger.Error("failed to save PostHistory for post: ", dbPost.ID, "; error: ", err)
			log.Printf("failed to save PostHistory for post %s: %v", post.ID, err)
			continue
		}
	}

	return nil
}

// Helper functions for parsing
func parsePrice(price string) int64 {
	price = strings.ReplaceAll(price, ",", "")
	price = strings.ReplaceAll(price, " تومان", "")
	price = strings.ReplaceAll(price, "ریال", "")
	price = replaceDigits(price) // اگر اعداد فارسی وجود دارند
	parsed, _ := strconv.ParseInt(price, 10, 64)
	return parsed
}

func parseArea(area string) int {
	area = replaceDigits(area)
	parsed, _ := strconv.Atoi(area)
	return parsed
}

func parseBedrooms(rooms string) int {
	rooms = replaceDigits(rooms)
	parsed, _ := strconv.Atoi(rooms)
	return parsed
}

func parseAge(yearBuilt string) uint8 {
	currentYear := time.Now().Year()
	yearBuilt = replaceDigits(yearBuilt)
	builtYear, _ := strconv.Atoi(yearBuilt)
	if builtYear > 0 {
		return uint8(currentYear - builtYear)
	}
	return 0
}

func parseFloors(floor string) uint8 {
	floor = replaceDigits(strings.Split(floor, " ")[0]) // اولین بخش طبقه
	parsed, _ := strconv.Atoi(floor)
	return uint8(parsed)
}

func containsFeature(features []string, feature string) bool {
	for _, f := range features {
		if strings.Contains(f, feature) {
			return true
		}
	}
	return false
}
