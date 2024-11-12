package models

import (
	"testing"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the models to create tables
	db.AutoMigrate(&models.FilterItem{}, &models.PostHistory{})
	return db, nil
}

func TestFilterItemRepository(t *testing.T) {
	datab, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	repo := db.NewFilterItemRepository(datab)

	t.Run("Create", func(t *testing.T) {
		filter := models.FilterItem{
			PriceMin:         100000,
			PriceMax:         500000,
			City:             "Tehran",
			BedroomsMin:      2,
			BedroomsMax:      4,
			CreatedDateStart: time.Now().AddDate(-1, 0, 0),
			CreatedDateEnd:   time.Now(),
		}

		savedFilter, err := repo.Create(filter)
		assert.NoError(t, err)
		assert.NotZero(t, savedFilter.ID)
		assert.Equal(t, "Tehran", savedFilter.City)
	})

	t.Run("FindByID", func(t *testing.T) {
		filter := models.FilterItem{
			PriceMin: 200000,
			PriceMax: 800000,
			City:     "Tehran",
		}

		savedFilter, _ := repo.Create(filter)
		foundFilter, err := repo.FindByID(savedFilter.ID)
		assert.NoError(t, err)
		assert.Equal(t, savedFilter.ID, foundFilter.ID)
		assert.Equal(t, "Tehran", foundFilter.City)
	})

	t.Run("FindAll", func(t *testing.T) {
		filters, err := repo.FindAll()
		assert.NoError(t, err)
		assert.True(t, len(filters) > 0)
	})

	t.Run("Update", func(t *testing.T) {
		filter := models.FilterItem{
			PriceMin: 300000,
			PriceMax: 600000,
			City:     "Tehran",
		}
		savedFilter, _ := repo.Create(filter)

		updatedData := models.FilterItem{
			City: "Mashhad",
		}
		updatedFilter, err := repo.Update(savedFilter.ID, updatedData)
		assert.NoError(t, err)
		assert.Equal(t, "Mashhad", updatedFilter.City)
	})

	t.Run("Delete", func(t *testing.T) {
		filter := models.FilterItem{
			PriceMin: 400000,
			PriceMax: 700000,
			City:     "Tehran",
		}
		savedFilter, _ := repo.Create(filter)

		err := repo.Delete(savedFilter.ID)
		assert.NoError(t, err)

		_, err = repo.FindByID(savedFilter.ID)
		assert.Error(t, err)
	})

	t.Run("SearchPostHistory", func(t *testing.T) {
		// Seed some PostHistory data
		post := models.PostHistory{
			PostID:     1,
			Price:      300000,
			City:       "Tehran",
			BedroomNum: 3,
			BuyMode:    "buy",
			Building:   "apartment",
		}
		datab.Create(&post)

		filter := models.FilterItem{
			PriceMin:    200000,
			PriceMax:    400000,
			City:        "Tehran",
			BedroomsMin: 2,
			BedroomsMax: 4,
		}

		results, err := repo.SearchPostHistory(filter)
		assert.NoError(t, err)
		assert.True(t, len(results) > 0)
		assert.Equal(t, "Tehran", results[0].City)
	})
}