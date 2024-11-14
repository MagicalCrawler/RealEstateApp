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

	"github.com/MagicalCrawler/RealEstateApp/models"
)

const (
	timeout = 10
)

func getOrCreateUserRunCommand(message *Message) models.User {
	// Check if user already exists by Telegram ID
	empty_user := models.User{}
	user, err := userRepository.FindByTelegramID(uint64(message.From.ID))
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		sendMessage(message.Chat.ID, "There was an error checking your profile. Please try again later.")
		return empty_user
	}
	if user == empty_user {
		// If the user does not exist, create a new user
		user = models.User{TelegramID: uint64(message.From.ID), Role: models.Role(models.USER)}
		log.Printf("User with id : %d created with role regular", message.From.ID)
		_, err = userRepository.Save(user)
		if err != nil {
			log.Printf("Error saving new user: %v", err)
			sendMessage(message.Chat.ID, "There was an error creating your profile. Please try again later.")
			return empty_user
		}
	}
	return user
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
func deleteMessage(chatID int, messageID int) error {
	data := url.Values{}
	data.Set("chat_id", strconv.Itoa(chatID))
	data.Set("message_id", strconv.Itoa(messageID))

	resp, err := http.PostForm(fmt.Sprintf("%s/deleteMessage", apiURL), data)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return err
	}
	defer resp.Body.Close()

	return nil
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
func getKeyboard(role models.Role) ReplyKeyboardMarkupWithLocation {
	switch {
	case role == models.ADMIN:
		return ReplyKeyboardMarkupWithLocation{
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
	case role == models.SUPER_ADMIN:
		return ReplyKeyboardMarkupWithLocation{
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
	case role == models.USER:
		return ReplyKeyboardMarkupWithLocation{
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
	default:
		return ReplyKeyboardMarkupWithLocation{
			Keyboard: [][]KeyboardButton{
				{
					{Text: "Help"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	}

}
func getWelcomeMessage(name string, role models.Role) string {
	switch {
	case role == models.USER:
		return fmt.Sprintf("Welcome %s!", name)
	case role == models.ADMIN:
		return fmt.Sprintf("Hi %s\nWelcome to admin panel !", name)
	case role == models.SUPER_ADMIN:
		return fmt.Sprintf("Hi %s\nWelcome to superadmin panel !", name)
	default:
		return "Welcome!"
	}
}

func sendHelpMessage(chatID int, text string) {

	sendMessageWithKeyboard(chatID, text, getKeyboard(models.USER))
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
					RequestLocation: true,
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	sendMessageWithKeyboard(chatID, "Please share your location:", keyboard)
}
