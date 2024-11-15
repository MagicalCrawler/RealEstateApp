package services

import (
	"encoding/json"
	"fmt"
	crawlerModels "github.com/MagicalCrawler/RealEstateApp/models/crawler"
	"github.com/MagicalCrawler/RealEstateApp/utils"
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
}

// NewCityService creates a new instance of CityService
func NewCityService() *CityService {
	return &CityService{
		cacheDuration: 30 * time.Minute,
	}
}

// GetCities returns the list of cities, using cache
func (s *CityService) GetCities() ([]crawlerModels.City, error) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if time.Since(s.lastUpdated) < s.cacheDuration && s.cache != nil {
		return s.cache, nil
	}

	resp, err := http.Get("https://api.divar.ir/v8/places/cities?level=all")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cityResponse crawlerModels.CityResponse
	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		return nil, fmt.Errorf("failed to decode city response: %w", err)
	}

	s.lastUpdated = time.Now()

	// Filter cities based on app settings
	provincialCenters, err := utils.LoadAppSettingsFile()
	if err != nil {
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
