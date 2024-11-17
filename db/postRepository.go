package db

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"gorm.io/gorm"
)

type PostRepo interface {
	PostIsExist(post models.Post) bool
	PostHistoryIsExist(postHistory models.PostHistory) bool
	FindByUnicode(UniCode string) (models.Post, models.PostHistory, error)
	PostSaving(uniCode string, src types.WebsiteSource) (models.Post, error)
	PostHistorySaving(postHistory models.PostHistory, post models.Post, crawlHistory models.CrawlHistory) (models.PostHistory, error)
	CrawlHistorySaving(crawlHistory models.CrawlHistory) (models.CrawlHistory, error)

	GetAllCrawlHistory() []models.CrawlHistory

	CrawlHistoryIsExist(crawlHistory models.CrawlHistory) bool
	GetMostVisitedPost() ([]models.PostHistory, error)
	GetAllPosts() ([]models.PostHistory, error)

}

type PostRepository struct {
	dbConnection *gorm.DB
	logger       *slog.Logger
}

func NewPostRepository(dbConnection *gorm.DB) PostRepo {
	return PostRepository{dbConnection: dbConnection}
}
func (pr PostRepository) CrawlHistoryIsExist(crawlHistory models.CrawlHistory) bool {
	var isExist bool
	err := pr.dbConnection.Table("crawl_histories").
		Select("count(*) > 0").
		Where("id = ?", crawlHistory.ID).
		Find(&isExist).Error
	if err != nil {
		return false
	}
	return isExist
}
func (pr PostRepository) GetMostVisitedPost() ([]models.PostHistory, error) {
	var posts []models.PostHistory

	err := pr.dbConnection.Table("posts").
		Select("posts.unique_code, posts.watched_num, post_histories.title, post_histories.post_url, post_histories.price").
		Joins("INNER JOIN post_histories ON post_histories.post_id = posts.id").
		Order("posts.watched_num DESC").
		Limit(10).
		Scan(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (pr PostRepository) GetAllPosts() ([]models.PostHistory, error) {
	var posts []models.PostHistory

	err := pr.dbConnection.Table("posts").
		Select("DISTINCT posts.unique_code, posts.watched_num, post_histories.title, post_histories.post_url, post_histories.price").
		Joins("INNER JOIN post_histories ON post_histories.post_id = posts.id").
		Order("posts.watched_num DESC").
		Scan(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (pr PostRepository) PostIsExist(post models.Post) bool {
	var isExist bool
	err := pr.dbConnection.Table("posts").Select("count(*) > 0").Where("unique_code = ?", post.UniqueCode).Find(&isExist).Error
	if err != nil {
		return false
	} else {
		return isExist
	}
}
func (pr PostRepository) PostSaving(uniCode string, src types.WebsiteSource) (models.Post, error) {
	post := models.Post{
		UniqueCode: uniCode,
		Website:    src,
	}
	if !pr.PostIsExist(post) {
		err := pr.dbConnection.Create(&post).Error
		return post, err
	}
	return post, errors.New("Post already exists")
}

func Separate(url string) (types.WebsiteSource, string) {
	var webSite types.WebsiteSource
	mySlice := strings.SplitN(url, "/", 6)
	if mySlice[2] == "divar.ir" {
		webSite = types.Divar
	}
	uniqueCode := mySlice[5]
	return webSite, uniqueCode
}

func (pr PostRepository) PostHistoryIsExist(postHistory models.PostHistory) bool {
	var isExist bool
	err := pr.dbConnection.Table("post_histories").Select("count(*) > 0").Where("post_url = ?", postHistory.PostURL).Find(&isExist).Error
	if err != nil {
		return false
	} else {
		return isExist
	}

}
func (pr PostRepository) FindByUnicode(UniCode string) (models.Post, models.PostHistory, error) {
	var post models.Post
	var postHistory models.PostHistory
	err := pr.dbConnection.First(&post, "unique_code = ?", UniCode).Error
	pr.dbConnection.Where("post_id = ?", post.ID).Find(&postHistory)
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
		HasParking:     postHistory.HasParking,
		ImageURL:       postHistory.ImageURL,
		Description:    postHistory.Description,
		CrawlHistory:   crawlHistory,
		CrawlHistoryID: crawlHistory.ID,
	}

	if !daba.PostHistoryIsExist(myPostHistory) {
		err := daba.dbConnection.Create(&myPostHistory).Error
		return myPostHistory, err
	}
	return myPostHistory, errors.New("Post history already exists")
}

func (dba PostRepository) CrawlHistorySaving(crawlHistory models.CrawlHistory) (models.CrawlHistory, error) {
	myCrawlHistory := models.CrawlHistory{
		PostNum:     crawlHistory.PostNum,
		CpuUsage:    crawlHistory.CpuUsage,
		MemoryUsage: crawlHistory.MemoryUsage,
		RequestsNum: crawlHistory.RequestsNum,
		StartedAt:   crawlHistory.StartedAt,
		FinishedAt:  crawlHistory.FinishedAt,
	}

	err := dba.dbConnection.Create(&myCrawlHistory).Error
	return myCrawlHistory, err
}

func (pr PostRepository) GetAllCrawlHistory() []models.CrawlHistory {
	var crawlHistories []models.CrawlHistory
	pr.dbConnection.Find(&crawlHistories)
	return crawlHistories
}
