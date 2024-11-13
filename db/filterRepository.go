package db

import (
	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type FilterItemRepository interface {
	Create(filterItem models.FilterItem) (models.FilterItem, error)
	FindByID(id uint) (models.FilterItem, error)
	FindAll() ([]models.FilterItem, error)
	Update(id uint, updatedData models.FilterItem) (models.FilterItem, error)
	Delete(id uint) error
	SearchPostHistory(filter models.FilterItem) ([]models.PostHistory, error)
}

type FilterItemRepositoryImpl struct {
	dbConnection *gorm.DB
}

func NewFilterItemRepository(dbConnection *gorm.DB) FilterItemRepository {
	return &FilterItemRepositoryImpl{dbConnection: dbConnection}
}


func (repo FilterItemRepositoryImpl) Create(filterItem models.FilterItem) (models.FilterItem, error) {
	err := repo.dbConnection.Create(&filterItem).Error
	return filterItem, err
}

func (repo FilterItemRepositoryImpl) FindByID(id uint) (models.FilterItem, error) {
	var filterItem models.FilterItem
	err := repo.dbConnection.First(&filterItem, id).Error
	return filterItem, err
}

func (repo FilterItemRepositoryImpl) FindAll() ([]models.FilterItem, error) {
	var filterItems []models.FilterItem
	err := repo.dbConnection.Find(&filterItems).Error
	return filterItems, err
}

func (repo FilterItemRepositoryImpl) Update(id uint, updatedData models.FilterItem) (models.FilterItem, error) {
	var filterItem models.FilterItem
	if err := repo.dbConnection.First(&filterItem, id).Error; err != nil {
		return filterItem, err
	}

	// Update fields based on updatedData
	err := repo.dbConnection.Model(&filterItem).Updates(updatedData).Error
	return filterItem, err
}

func (repo FilterItemRepositoryImpl) Delete(id uint) error {
	if err := repo.dbConnection.Delete(&models.FilterItem{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo FilterItemRepositoryImpl) SearchPostHistory(filter models.FilterItem) ([]models.PostHistory, error) {
	var posts []models.PostHistory
	query := repo.dbConnection.Model(&models.PostHistory{})

	// Apply filter conditions
	if filter.PriceMin > 0 {
		query = query.Where("price >= ?", filter.PriceMin)
	}
	if filter.PriceMax > 0 {
		query = query.Where("price <= ?", filter.PriceMax)
	}
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.Neighborhood != "" {
		query = query.Where("neighbor = ?", filter.Neighborhood)
	}
	if filter.AreaMin > 0 {
		query = query.Where("area >= ?", filter.AreaMin)
	}
	if filter.AreaMax > 0 {
		query = query.Where("area <= ?", filter.AreaMax)
	}
	if filter.BedroomsMin > 0 {
		query = query.Where("bedroom_num >= ?", filter.BedroomsMin)
	}
	if filter.BedroomsMax > 0 {
		query = query.Where("bedroom_num <= ?", filter.BedroomsMax)
	}
	if filter.Category != "" {
		query = query.Where("buy_mode = ?", filter.Category)
	}
	if filter.AgeMin > 0 {
		query = query.Where("age >= ?", filter.AgeMin)
	}
	if filter.AgeMax > 0 {
		query = query.Where("age <= ?", filter.AgeMax)
	}
	if filter.PropertyType != "" {
		query = query.Where("building = ?", filter.PropertyType)
	}
	if filter.FloorMin > 0 {
		query = query.Where("floors_num >= ?", filter.FloorMin)
	}
	if filter.FloorMax > 0 {
		query = query.Where("floors_num <= ?", filter.FloorMax)
	}
	if filter.HasStorage {
		query = query.Where("has_ware_house = ?", filter.HasStorage)
	}
	if filter.HasElevator {
		query = query.Where("has_elevator = ?", filter.HasElevator)
	}
	if !filter.CreatedDateStart.IsZero() {
		query = query.Where("created_at >= ?", filter.CreatedDateStart)
	}
	if !filter.CreatedDateEnd.IsZero() {
		query = query.Where("created_at <= ?", filter.CreatedDateEnd)
	}

	err := query.Find(&posts).Error
	return posts, err
}
