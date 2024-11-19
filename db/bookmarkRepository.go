package db

import (
	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

type BookmarkRepo interface {
	Find(bookmark models.Bookmark) (models.Bookmark, error)
	FindAll(UserID uint) ([]models.Bookmark, error)
	Save(post models.Post, user models.User) error
	Delete(bookmark models.Bookmark) error
}

// Connection to database
type BookmarkRepositoryImpl struct {
	dbConnection *gorm.DB
}

// connection in interface
func NewBookmarkRepository(dbConnection *gorm.DB) BookmarkRepo {
	return BookmarkRepositoryImpl{dbConnection: dbConnection}
}

// find a bookmark by an isntance of its model
func (br BookmarkRepositoryImpl) Find(bookmark models.Bookmark) (models.Bookmark, error) {
	if err := br.dbConnection.Find(&bookmark).Error; err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

// find all bookmark of a user
func (br BookmarkRepositoryImpl) FindAll(userID uint) ([]models.Bookmark, error) {
	bookmarks := []models.Bookmark{}
	if err := br.dbConnection.Find(&bookmarks).Error; err != nil {
		return []models.Bookmark{}, err
	}
	return bookmarks, nil

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
