package models

import (
	"fmt"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/types"
	"log"
	"testing"
)

func TestFindBookmark(t *testing.T) {
	defer clearData()
	user := models.User{
		TelegramID: 125478,
		Role:       models.USER,
	}
	dbConnection.Create(&user)
	post := models.Post{
		UniqueCode: "aabbcc",
		Website:    types.Divar,
	}
	dbConnection.Create(&post)
	bookmark := models.Bookmark{
		Post:   post,
		PostID: post.ID,
		User:   user,
		UserID: user.ID,
	}
	bookmarkRep := db.NewBookmarkRepository(dbConnection)
	bookmarkRep.Save(post, user)
	if myBookMark, err := bookmarkRep.Find(bookmark); err != nil {
		log.Fatal(`bookmark already exists`)
	} else {
		fmt.Println(myBookMark)
	}
}

func TestDelete(t *testing.T) {
	//defer clearData()
	user := models.User{
		TelegramID: 125478,
		Role:       models.USER,
	}
	dbConnection.Create(&user)
	post := models.Post{
		UniqueCode: "aabbcc",
		Website:    types.Divar,
	}
	dbConnection.Create(&post)
	post1 := models.Post{
		UniqueCode: "ddccee",
		Website:    types.Divar,
	}
	dbConnection.Create(&post1)
	bookmarkRep := db.NewBookmarkRepository(dbConnection)
	bookmarkRep.Save(post, user)
	bookmarkRep.Save(post1, user)
	bookmark := models.Bookmark{
		Post:   post,
		PostID: post.ID,
		User:   user,
		UserID: user.ID,
	}
	fmt.Println(bookmarkRep.Delete(bookmark))
}
