package client

import (
	"fmt"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"gorm.io/gorm"
)

var userRepository db.UserRepository
var apiURL string

func Run(dbConnection *gorm.DB) {
	userRepository = db.CreateNewUserRepository(dbConnection)
	apiURL = "https://api.telegram.org/bot" + utils.GetConfig("TELEGRAM_TOKEN")
	go pollUpdates()

	fmt.Println("Bot is running...")
	select {}
}

func handleMessage(message *Message) {
	deleteMessage(message.Chat.ID, message.MessageID-1)
	deleteMessage(message.Chat.ID, message.MessageID-2)

	user := getOrCreateUserRunCommand(message)

	switch {
	case message.Text == "/start":
		sendMessageWithKeyboard(message.Chat.ID, getWelcomeMessage(message.From.FirstName, user.Role), getKeyboard(user.Role))
		return

	case message.Text == "Help":
		msgHelp := `Real Estate Finder Bot!
					/search to find properties based on filters like price, location, and type.
					/notify to get alerts for new listings matching your preferences.
					/help for more information.`
		sendMessageWithKeyboard(message.Chat.ID, msgHelp, getKeyboard(user.Role))
		return
	case message.Text == "Send Location":
		msg := "You can send me location with your telegram attachment"
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Location.Latitude != 0:
		msg := fmt.Sprintf("Your selected location is with latitude: %d, and longitude%d", message.Location.Latitude, message.Location.Longitude)
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Text == "Filters":
		msg := "You entered filters"
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	default:
		msg := "I didn't understand that command."
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	}
}
