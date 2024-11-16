package services

import (
	"encoding/json"
	"fmt"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// CityService handles fetching and caching city data
type CityService struct {
	cache         []crawlerModels.City
	cacheDuration time.Duration
	cacheMutex    sync.Mutex
	lastUpdated   time.Time
	logger        *slog.Logger
}

// NewCityService creates a new instance of CityService
func NewCityService() *CityService {
	return &CityService{
		cacheDuration: 6 * time.Hour,
		logger:        utils.NewLogger("City_Service"),
	}
}

// GetCities returns the list of cities, using cache
func (s *CityService) GetCities() ([]crawlerModels.City, error) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if time.Since(s.lastUpdated) < s.cacheDuration && s.cache != nil {
		return s.cache, nil
	}

	divarAllCityAPIAddress := utils.GetConfig("API_CITIES_URL")
	resp, err := http.Get(divarAllCityAPIAddress)
	if err != nil {
		s.logger.Error("failed to fetch cities ", err)
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexpected status code ", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cityResponse crawlerModels.CityResponse
	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		s.logger.Error("failed to decode response ", err)
		return nil, fmt.Errorf("failed to decode city response: %w", err)
	}

	s.lastUpdated = time.Now()

	// Filter cities based on app settings
	provincialCenters, err := utils.LoadAppSettingsFile()
	if err != nil {
		s.logger.Error("failed to load app settings ", err)
		return nil, fmt.Errorf("failed to load app settings: %w", err)
	}

	var filteredCities []crawlerModels.City
	for _, city := range cityResponse.Cities {
		for _, center := range provincialCenters {
			if city.Name == center.Name {
				filteredCities = append(filteredCities, city)
				break
			}
		}
	}

	s.cache = filteredCities
	return s.cache, nil
}
