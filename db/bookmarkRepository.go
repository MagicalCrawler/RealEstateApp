package db

import (
	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type BookmarkRepo interface {
	Find(bookmark models.Bookmark) (models.Bookmark, error)
	Save(post models.Post, user models.User) error
	Delete(bookmark models.Bookmark) error
}

type BookmarkRepositoryImpl struct {
	dbConnection *gorm.DB
}

func NewBookmarkRepository(dbConnection *gorm.DB) BookmarkRepo {
	return BookmarkRepositoryImpl{dbConnection: dbConnection}
}

func (br BookmarkRepositoryImpl) Find(bookmark models.Bookmark) (models.Bookmark, error) {
	if err := br.dbConnection.Find(&bookmark).Error; err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil

}

func (br BookmarkRepositoryImpl) Save(post models.Post, user models.User) error {
	bookmark := models.Bookmark{
		Post:   post,
		PostID: post.ID,
		User:   user,
		UserID: user.ID,
	}
	return br.dbConnection.Create(&bookmark).Error
}

func (br BookmarkRepositoryImpl) Delete(bookmark models.Bookmark) error {
	return br.dbConnection.Where("post_id = ? AND user_id = ?", bookmark.PostID, bookmark.UserID).Delete(&models.Bookmark{}).Error
}