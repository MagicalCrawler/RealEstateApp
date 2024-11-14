package client

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
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
		msg := "You can send me location with your telegram attachment.."
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Location.Latitude != 0:
		msg := fmt.Sprintf("Your selected location is with latitude: %f, and longitude%f👌\n\nNow send me your desired radius with pattern👉 \"rediuse=<number>\"", message.Location.Latitude, message.Location.Longitude)

		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case strings.Contains(message.Text, "rediuse="):
		msg := fmt.Sprintf("You entered radius: %s", message.Text[8:])
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Text == "Filters":
		msg := "You entered filters"
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Text == "Premium":
		msg := "Send me user id with pattern👉 \"Id=<number>\""
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case strings.Contains(message.Text, "Id="):
		msg := fmt.Sprintf("User with id :%s changed to Premium client.", message.Text)
		id, err := strconv.ParseInt(message.Text[3:], 10, 64)
		if err != nil {
			msg = "Invalid ID format. Please use 'Id=<number>'."
		} else {
			_, err := userRepository.UpdateUserType(uint(id), models.PREMIUM)
			if err != nil {
				msg = fmt.Sprintf("Error updating user type: %v", err)
			}
		}
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Text == "Clients":
		msg := "All Clients:\n"
		users, err := userRepository.FindAllUsersByRole(models.USER)
		if err != nil {
			msg = "Error fetching clients: "
			log.Fatal(msg + err.Error())
		} else {
			if len(users) == 0 {
				msg = "No clients found"
			} else {
				for _, u := range users {
					msg += fmt.Sprintf("   ID: %d, TelegramID: %d\n", u.ID, u.TelegramID)
					msg += fmt.Sprint("-------------------------\n")
				}
			}
		}
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	case message.Text == "Admins":
		msg := "All Admins:\n"
		users, err := userRepository.FindAllUsersByRole(models.ADMIN)
		if err != nil {
			msg = "Error fetching clients: "
			log.Fatal(msg + err.Error())
		} else {
			if len(users) == 0 {
				msg = "No admins found"
			} else {
				for _, u := range users {
					msg += fmt.Sprintf("   ID: %d, TelegramID: %d\n", u.ID, u.TelegramID)
					msg += fmt.Sprint("-------------------------\n")
				}
			}
		}
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	default:
		msg := "I didn't understand that command."
		sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
		return
	}
}
