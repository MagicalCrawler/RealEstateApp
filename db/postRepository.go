package db

import (
	"errors"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
	"strings"
)

type PostRepo interface {
	PostIsExist(dbConnection *gorm.DB, post models.Post) bool
	PostHistoryIsExist(dbConnection *gorm.DB, postHistory models.PostHistory) bool
	Find(UniCode string) (models.Post, models.PostHistory, error)
	PostHistorySaving(postHistory models.PostHistory, post models.Post, crawlHistory models.CrawlHistory) (models.PostHistory, error)
	PostSaving(uniCode string, src types.WebsiteSource) (models.Post, error)
	separate(url string) (types.WebsiteSource, string)
}

type PostRepository struct {
	dbConnection *gorm.DB
}

func NewPostRepository(dbConnection *gorm.DB) PostRepository {
	return PostRepository{dbConnection: dbConnection}
}

func PostIsExist(dbConnection *gorm.DB, post models.Post) bool {
	var isExist bool
	err := dbConnection.Table("posts").Select("count(*) > 0").Where("unique_code = ?", post.UniqueCode).Find(&isExist).Error
	if err != nil {
		return false
	} else {
		return isExist
	}
	return false
}
func (daba PostRepository) PostSaving(uniCode string, src types.WebsiteSource) (models.Post, error) {
	post := models.Post{
		UniqueCode: uniCode,
		Website:    src,
	}
	if !PostIsExist(daba.dbConnection, post) {
		err := daba.dbConnection.Create(&post).Error
		return post, err
	}
	return post, errors.New("Post already exists")
}
func separate(url string) (types.WebsiteSource, string) {
	var webSite types.WebsiteSource
	mySlice := strings.SplitN(url, "/", 6)
	if mySlice[2] == "divar.ir" {
		webSite = types.Divar
	}
	uniqueCode := mySlice[5]
	return webSite, uniqueCode
}

func PostHistoryIsExist(dbConnection *gorm.DB, postHistory models.PostHistory) bool {
	var isExist bool
	err := dbConnection.Table("post_histories").Select("count(*) > 0").Where("post_url = ?", postHistory.PostURL).Find(&isExist).Error
	if err != nil {
		return false
	} else {
		return isExist
	}
	return false
}
func (daba PostRepository) Find(UniCode string) (models.Post, models.PostHistory, error) {
	var post models.Post
	var postHistory models.PostHistory
	err := daba.dbConnection.First(&post, "unique_code = ?", UniCode).Error
	daba.dbConnection.Where("post_id = ?", post.ID).Find(&postHistory)
	return post, postHistory, err

}
func (daba PostRepository) PostHistorySaving(postHistory models.PostHistory, post models.Post, crawlHistory models.CrawlHistory) (models.PostHistory, error) {

	myPostHistory := models.PostHistory{
		Post:           post,
		PostID:         post.ID,
		Title:          postHistory.Title,
		PostURL:        postHistory.PostURL,
		Price:          postHistory.Price,
		Deposit:        postHistory.Deposit,
		Rent:           postHistory.Rent,
		City:           postHistory.City,
		Neighborhood:   postHistory.Neighborhood,
		Area:           postHistory.Area,
		BedroomNum:     postHistory.BedroomNum,
		BuyMode:        postHistory.BuyMode,
		Building:       postHistory.Building,
		Age:            postHistory.Age,
		FloorsNum:      postHistory.FloorsNum,
		HasStorage:     postHistory.HasStorage,
		HasElevator:    postHistory.HasElevator,
		HsaParking:     postHistory.HsaParking,
		ImageURL:       postHistory.ImageURL,
		Description:    postHistory.Description,
		CrawlHistory:   crawlHistory,
		CrawlHistoryID: crawlHistory.ID,
	}
	if !PostHistoryIsExist(daba.dbConnection, myPostHistory) {
		err := daba.dbConnection.Create(&myPostHistory).Error
		return myPostHistory, err
	}
	return myPostHistory, errors.New("Post history already exists")
}
