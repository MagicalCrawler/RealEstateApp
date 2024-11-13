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
	msgHelp = `Real Estate Finder Bot!
	/search to find properties based on filters like price, location, and type.
	/notify to get alerts for new listings matching your preferences.
	/help for more information.`
	timeout = 10
)

var (
	apiURL         string
	userRepository db.UserRepository
	firstMenu      = []string{
		"Search", "Setting", "Populars", "Help",
	}
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
				handleCallback(update.Callback)
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
			if existingUser.Role == 2 {
				log.Printf("User with id : %d enters with role regular\n", existingUser.TelegramID)
			} else if existingUser.Role == 1 {
				log.Printf("User with id : %d enters with role admin\n", existingUser.TelegramID)
			} else if existingUser.Role == 0 {
				log.Printf("User with id : %d enters with role superadmin\n", existingUser.TelegramID)
			}
			// User already exists, so just show a welcome message
			sendMainMenu(message.Chat.ID, message.From.FirstName)
			return
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

		sendMainMenu(message.Chat.ID, message.From.FirstName)
		return

	case message.Text == "Choose an option":
		sendMainMenuSelectionInlineKeyboard(message.Chat.ID, "Menu item:")
	default:
		sendMessage(message.Chat.ID, "I didn't understand that command.")
	}
}

func sendMainMenu(chatID int, name string) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Choose an option"},
				{Text: "Share Location", RequestLocation: true},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	welcomeMsg := fmt.Sprintf("Welcome %s!", name)
	sendMessageWithKeyboard(chatID, welcomeMsg, keyboard)
}

func sendMainMenuSelectionInlineKeyboard(chatID int, text string) {
	keyboard := createMainMenuInlineKeyboard()
	sendMessageWithInlineKeyboard(chatID, text, keyboard)
}

func createMainMenuInlineKeyboard() InlineKeyboardMarkup {
	rows := [][]InlineKeyboardButton{}
	row := []InlineKeyboardButton{}

	for i, option := range firstMenu {
		row = append(row, InlineKeyboardButton{
			Text: option,
			Data: option,
		})

		if (i+1)%4 == 0 {
			rows = append(rows, row)
			row = []InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	return InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func handleCallback(callback *CallbackQuery) {
	selectedOption := callback.Data
	if selectedOption == "Help" {
		sendHelpMessage(callback.Message.Chat.ID, msgHelp)
	}
	answerCallbackQuery(callback.ID, "Selected!")
}

func sendHelpMessage(chatID int, text string) {
	keyboard := createMainMenuInlineKeyboard()
	sendMessageWithInlineKeyboard(chatID, text, keyboard)
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
