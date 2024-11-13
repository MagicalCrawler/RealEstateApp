package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/models"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"gorm.io/gorm"
)

const (
	timeout = 10
)

var (
	apiURL         string
	userRepository db.UserRepository
)

func InitUserRepository(dbConnection *gorm.DB) {
	userRepository = db.CreateNewUserRepository(dbConnection)
}

func Run(dbConnection *gorm.DB) {
	InitUserRepository(dbConnection)

	apiURL = "https://api.telegram.org/bot" + utils.GetConfig("TELEGRAM_TOKEN")
	go pollUpdates()

	fmt.Println("Bot is running...")
	select {}
}
func pollUpdates() {
	offset := 0

	for {
		updates, err := getUpdates(offset)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			continue
		}

		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message != nil {
				handleMessage(update.Message)
			} else if update.Callback != nil {
				// handleCallback(update.Callback)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(fmt.Sprintf("%s/getUpdates?offset=%d&timeout=%d", apiURL, offset, timeout))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func handleMessage(message *Message) {
	switch {
	case message.Text == "/start":

		// Check if user already exists by Telegram ID
		existingUser, err := userRepository.FindByTelegramID(uint64(message.From.ID))
		if err != nil {
			log.Printf("Error checking if user exists: %v", err)
			sendMessage(message.Chat.ID, "There was an error checking your profile. Please try again later.")
			return
		}
		empty_user := models.User{}
		if *existingUser != empty_user {
			// User already exists, so just show a welcome message
			if existingUser.Role == 2 {
				log.Printf("User with id : %d enters with role regular\n", existingUser.TelegramID)
				sendMenuForRegularUser(message.Chat.ID, message.From.FirstName)
				return
			} else if existingUser.Role == 1 {
				log.Printf("User with id : %d enters with role admin\n", existingUser.TelegramID)
				sendMenuForAdminUser(message.Chat.ID, message.From.FirstName)
				return
			} else if existingUser.Role == 0 {
				log.Printf("User with id : %d enters with role superadmin\n", existingUser.TelegramID)
				sendMenuForSuperAdminUser(message.Chat.ID, message.From.FirstName)
				return
			}
		}

		// If the user does not exist, create a new user
		newUser := models.User{TelegramID: uint64(message.From.ID), Role: models.Role(models.USER)}
		log.Println("User with id : %d created with role regular", message.From.ID)
		_, err = userRepository.Save(newUser)
		if err != nil {
			log.Printf("Error saving new user: %v", err)
			sendMessage(message.Chat.ID, "There was an error creating your profile. Please try again later.")
			return
		}

		sendMenuForRegularUser(message.Chat.ID, message.From.FirstName)
		return
	case message.Text == "Help":
		msgHelp := `Real Estate Finder Bot!
					/search to find properties based on filters like price, location, and type.
					/notify to get alerts for new listings matching your preferences.
					/help for more information.`
		sendHelpMessage(message.Chat.ID, msgHelp)
		return
	case message.Text == "Send Location":
		msg := "You can send me location with your telegram attachment"
		sendMessage(message.Chat.ID, msg)
		return
	case message.Location.Latitude != 0:
		msg := fmt.Sprintf("Your selected location is with latitude: %d, and longitude%d", message.Location.Latitude, message.Location.Longitude)
		sendMessage(message.Chat.ID, msg)
		return
	case message.Text == "Filters":
	default:

		sendMessage(message.Chat.ID, "I didn't understand that command.")
		return
	}
}

// func sendMenuForRegularUser(chatID int, name string) {
// 	keyboard := ReplyKeyboardMarkupWithLocation{
// 		Keyboard: [][]KeyboardButton{
// 			{
// 				{Text: "Choose an option"},
// 				{Text: "Share Location", RequestLocation: true},
// 			},
// 		},
// 		ResizeKeyboard:  true,
// 		OneTimeKeyboard: true,
// 	}
// 	welcomeMsg := fmt.Sprintf("Welcome %s!", name)
// 	sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
// }

func sendMenuForAdminUser(chatID int, name string) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Filters"},
				{Text: "Premium"},
				{Text: "Errors"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	welcomeMsg := fmt.Sprintf("Welcome to admin panel %s!", name)
	sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
}
func sendMenuForSuperAdminUser(chatID int, name string) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Admins"},
				{Text: "Premium"},
				{Text: "Clients"},
			},
			{
				{Text: "Monitor"},
				{Text: "ََAdvertisements"},
			},
			{
				{Text: "Crawler Setting"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	welcomeMsg := fmt.Sprintf("Welcome to super admin panel %s!", name)
	sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
}
func sendMenuForRegularUser(chatID int, name string) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Search"},
				{Text: "Setting"},
				{Text: "Populars"},
			},
			{
				{Text: "Send Location"},
				{Text: "Help"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	welcomeMsg := fmt.Sprintf("Welcome %s!", name)
	// keyboard := createMainMenuInlineKeyboard()
	sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
}

// func handleCallback(callback *CallbackQuery) {
// 	selectedOption := callback.Data
// 	if selectedOption == "Help" {
// 		msgHelp := `Real Estate Finder Bot!
// 					/search to find properties based on filters like price, location, and type.
// 					/notify to get alerts for new listings matching your preferences.
// 					/help for more information.`
// 		sendHelpMessage(callback.Message.Chat.ID, msgHelp)
// 	}
// 	answerCallbackQuery(callback.ID, "Selected!")
// }

func sendHelpMessage(chatID int, text string) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Search"},
				{Text: "Setting"},
				{Text: "Populars"},
			},
			{
				{Text: "Send Location"},
				{Text: "Help"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	sendMessageWithKeyboard(chatID, text, keyboard)
}

func answerCallbackQuery(callbackID, text string) {
	payload := map[string]interface{}{
		"callback_query_id": callbackID,
		"text":              text,
		"show_alert":        false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	_, err = http.Post(fmt.Sprintf("%s/answerCallbackQuery", apiURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error answering callback query: %v", err)
	}
}

func sendMessageWithInlineKeyboard(chatID int, text string, keyboard InlineKeyboardMarkup) {
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/sendMessage", apiURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error sending message with inline keyboard: %v", err)
		return
	}
	defer resp.Body.Close()
}

func sendMessageWithKeyboard(chatID int, text string, keyboard ReplyKeyboardMarkupWithLocation) {
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/sendMessage", apiURL), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error sending message with keyboard: %v", err)
		return
	}
	defer resp.Body.Close()
}

func sendMessage(chatID int, text string) {
	data := url.Values{}
	data.Set("chat_id", strconv.Itoa(chatID))
	data.Set("text", text)

	resp, err := http.PostForm(fmt.Sprintf("%s/sendMessage", apiURL), data)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}
	defer resp.Body.Close()
}
func sendLocationRequest(chatID int) {
	// Set up a keyboard with a location request button
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				KeyboardButton{
					Text:            "Send your location 📍",
					RequestLocation: true, // Request location on button click
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	// Send message with location request keyboard
	sendMessageWithKeyboard(chatID, "Please share your location:", keyboard)
}
