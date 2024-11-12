package models

import (
	"testing"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWatchListTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the models to create tables
	db.AutoMigrate(&models.FilterItem{}, &models.PostHistory{}, &models.WatchList{})
	return db, nil
}

func TestWatchListRepository(t *testing.T) {
	// Setup test database
	datab, err := setupWatchListTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Create repository instance
	repo := db.NewWatchListRepository(datab)

	t.Run("Create", func(t *testing.T) {
		watchList := models.WatchList{
			UserID:          1,
			FilterItemID:    1,
			RefreshInterval: 10,
			LastChecked:     time.Now(),
		}

		// Test creating a WatchList
		savedWatchList, err := repo.Create(watchList)
		assert.NoError(t, err)
		assert.NotZero(t, savedWatchList.ID)
		assert.Equal(t, watchList.UserID, savedWatchList.UserID)
		assert.Equal(t, watchList.FilterItemID, savedWatchList.FilterItemID)
	})

	t.Run("FindByID", func(t *testing.T) {
		watchList := models.WatchList{
			UserID:          2,
			FilterItemID:    2,
			RefreshInterval: 15,
			LastChecked:     time.Now(),
		}
		// Save a WatchList for testing FindByID
		savedWatchList, _ := repo.Create(watchList)

		// Test finding a WatchList by ID
		foundWatchList, err := repo.FindByID(savedWatchList.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedWatchList.ID, foundWatchList.ID)
		assert.Equal(t, savedWatchList.UserID, foundWatchList.UserID)
		assert.Equal(t, savedWatchList.FilterItemID, foundWatchList.FilterItemID)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Seed the database with multiple WatchLists
		repo.Create(models.WatchList{
			UserID:          3,
			FilterItemID:    3,
			RefreshInterval: 5,
			LastChecked:     time.Now(),
		})
		repo.Create(models.WatchList{
			UserID:          4,
			FilterItemID:    4,
			RefreshInterval: 7,
			LastChecked:     time.Now(),
		})

		// Test finding all WatchLists
		watchLists, err := repo.FindAll()
		assert.NoError(t, err)
		assert.True(t, len(watchLists) > 0)
	})

	t.Run("Update", func(t *testing.T) {
		watchList := models.WatchList{
			UserID:          5,
			FilterItemID:    5,
			RefreshInterval: 20,
			LastChecked:     time.Now(),
		}
		// Save a WatchList for testing Update
		savedWatchList, _ := repo.Create(watchList)

		// Prepare updated data
		updatedData := models.WatchList{
			UserID:          6,
			FilterItemID:    6,
			RefreshInterval: 25,
		}

		// Test updating a WatchList
		updatedWatchList, err := repo.Update(savedWatchList.ID, updatedData)
		assert.NoError(t, err)
		assert.Equal(t, updatedData.UserID, updatedWatchList.UserID)
		assert.Equal(t, updatedData.FilterItemID, updatedWatchList.FilterItemID)
		assert.Equal(t, updatedData.RefreshInterval, updatedWatchList.RefreshInterval)
	})

	t.Run("Delete", func(t *testing.T) {
		watchList := models.WatchList{
			UserID:          7,
			FilterItemID:    7,
			RefreshInterval: 30,
			LastChecked:     time.Now(),
		}
		// Save a WatchList for testing Delete
		savedWatchList, _ := repo.Create(watchList)

		// Test deleting a WatchList
		err := repo.Delete(savedWatchList.ID)
		assert.NoError(t, err)

		// Test finding the deleted WatchList by ID
		_, err = repo.FindByID(savedWatchList.ID)
		assert.Error(t, err)
	})
}
