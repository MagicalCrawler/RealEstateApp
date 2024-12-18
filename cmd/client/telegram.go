package client

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

type Command interface {
	Execute(message *Message, user *models.User)
	AllowedRoles() []models.Role // returns roles that can execute this command
}

var (
	CommandRegistry    map[string]Command
	userRepository     db.UserRepository
	postRepository     db.PostRepo
	bookmarkRepository db.BookmarkRepo
	filterRepository   db.FilterItemRepository
	apiURL             string
)

func Run(userRepo db.UserRepository, postRepo db.PostRepo, bookmarkRepo db.BookmarkRepo, filterRepo db.FilterItemRepository) {
	postRepository = postRepo
	userRepository = userRepo
	bookmarkRepository = bookmarkRepo
	filterRepository = filterRepo

	apiURL = "https://api.telegram.org/bot" + utils.GetConfig("TELEGRAM_TOKEN")
	initializeCommands()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pollUpdates(ctx)

	log.Println("Bot is running...")
	select {}
}
func isRoleAllowed(userRole models.Role, allowedRoles []models.Role) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
func timedGoroutine() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// Cleanup logic when the context is canceled
				log.Println("Goroutine finished after timeout.")
				return
			default:
				// Simulate some work
				time.Sleep(1 * time.Second)
				log.Println("Goroutine is running.")
			}
		}
	}(ctx)
}
func handleMessage(message *Message) {
	deleteMessage(message.Chat.ID, message.MessageID-1)
	deleteMessage(message.Chat.ID, message.MessageID-2)

	user := getOrCreateUserRunCommand(message)
	if message.Location.Latitude != 0 {
		message.Title = "Location Attachment"
	} else if strings.Contains(message.Title, "Id=") {
		message.Value = message.Title
		message.Title = "Change To Premium"
	} else if strings.Contains(message.Title, "admin=") {
		message.Value = message.Title
		message.Title = "Create Admin"
	} else if strings.Contains(message.Title, "redius=") {
		message.Value = message.Title
		message.Title = "Get Redius"
	} else if strings.Contains(message.Title, "B=") {
		message.Value = message.Title
		message.Title = "Get Bookmark Id"
	} else if message.Title == "s" || message.Title == "d" {
		if message.Title == "s" {
			message.Value = "sheypoor"
		} else {
			message.Value = "divar"
		}
		message.Title = "Get Website"
	} else if message.Title == "c" {
		message.Title = "Get Admin Id"
	}
	if cmd, exists := CommandRegistry[message.Title]; exists {
		if isRoleAllowed(user.Role, cmd.AllowedRoles()) {
			cmd.Execute(message, &user)
			return
		} else {
			sendMessageWithKeyboard(message.Chat.ID, "You do not have permission to use this command.", getKeyboard(user.Role))
		}
	} else {
		saveUserFilterInput(message.Chat.ID, user.ID, message.Title)
		// sendMessageWithKeyboard(message.Chat.ID, "I didn't understand that command.", getKeyboard(user.Role))
	}
}
