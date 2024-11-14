package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/crawlers"
)

var (
	cityCache       []crawlers.City
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
	cacheDuration   = 30 * time.Minute // مدت زمان اعتبار داده‌های کش
)

// fetchCitiesFromAPIWithCache با استفاده از لایه کش داده‌ها را بازیابی می‌کند
func fetchCitiesFromAPIWithCache() ([]crawlers.City, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// چک کردن اعتبار کش
	if time.Now().Before(cacheExpiration) && cityCache != nil {
		return cityCache, nil
	}

	// در صورت نبودن داده در کش یا منقضی شدن، درخواست به API
	resp, err := http.Get("https://api.divar.ir/v8/places/cities?level=all")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cityResponse crawlers.CityResponse
	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		return nil, fmt.Errorf("failed to decode city response: %w", err)
	}

	// به‌روزرسانی داده‌ها و تنظیم اعتبار کش
	cityCache = cityResponse.Cities
	cacheExpiration = time.Now().Add(cacheDuration)

	return cityCache, nil
}
