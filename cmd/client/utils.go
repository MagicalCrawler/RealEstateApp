package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/models"
)

const (
	timeout = 10
)

var userLastFilterMap = make(map[uint]uint) // Stores the last filter selected by each user (UserID -> FilterID)
// Temporary in-memory storage for user filter items
var userFilterItems = make(map[uint]*models.FilterItem)

func getOrCreateUserRunCommand(message *Message) models.User {
	// Check if user already exists by Telegram ID
	empty_user := models.User{}
	user, err := userRepository.FindByTelegramID(uint64(message.From.ID))
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		sendMessage(message.Chat.ID, "There was an error checking your profile. Please try again later.")
		return empty_user
	}
	if user.ID == 0 {
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

func pollUpdates(ctx context.Context) {
	offset := 0

	for {
		select {
		case <-ctx.Done():
			// Context canceled or timed out, clean up and exit
			log.Println("pollUpdates: Stopping due to context cancellation.")
			return
		default:
			updates, err := getUpdates(offset)
			if err != nil {
				log.Printf("Error getting updates: %v", err)
				// Short delay before retrying to prevent tight error loop
				time.Sleep(1 * time.Second)
				continue
			}

			for _, update := range updates {
				offset = update.UpdateID + 1
				if update.Message != nil {
					handleMessage(update.Message)
				} else if update.Callback != nil {
					// handleCallback(update.Callback)
					handleCallbackQuery(update.Callback)
				}
			}

			// Avoid excessive API polling; sleep for 1 second between calls
			time.Sleep(1 * time.Second)
		}
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
					{Text: "Clients"},
				}, {
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
					{Text: "Monitor"},
					{Text: "Clients"},
				},
				{

					{Text: "Advertisements"},
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
					{Text: "Filter"},
					{Text: "Search"},
					{Text: "Bookmark"},
				},
				{
					// {Text: "Help"},
					{Text: "Send Location"},
					{Text: "Export CSV"},
				},
				{
					{Text: "Setting"},
					{Text: "Populars"},
				},
				{
					{Text: "Select Resource Website"},
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

func createInlineKeyboardFromOptions(options []string) InlineKeyboardMarkup {
	buttons := make([][]InlineKeyboardButton, 0)
	for _, option := range options {
		row := []InlineKeyboardButton{
			{
				Text: option,
				Data: option, // Use the filter name as callback data
			},
		}
		buttons = append(buttons, row)
	}
	return InlineKeyboardMarkup{InlineKeyboard: buttons}
}

// Temporary storage for user filters
var userFilters = make(map[uint64]map[string]string)

// Function to handle callback queries (filter selection)
func handleCallbackQuery(callbackQuery *CallbackQuery) {
	userID := uint64(callbackQuery.From.ID)
	_, err := userRepository.FindByTelegramID(userID)

	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return
	}

	selectedFilter := callbackQuery.Data
	chatID := int64(callbackQuery.Message.Chat.ID)

	// Initialize the user's filter map if it doesn't exist
	if _, exists := userFilters[userID]; !exists {
		userFilters[userID] = make(map[string]string)
	}

	// Prompt user for input based on the selected filter
	switch selectedFilter {
	case "Price Range":
		promptUserForInput(chatID, "Enter price range (e.g., 100000-200000):")
	case "City":
		promptUserForInput(chatID, "Enter city name:")
	case "Neighborhood":
		promptUserForInput(chatID, "Enter neighborhood name:")
	case "Area Range":
		promptUserForInput(chatID, "Enter area range (e.g., 50-200 square meters):")
	case "Bedroom Count Range":
		promptUserForInput(chatID, "Enter bedroom count range (e.g., 1-3):")
	case "Category (Rent/Buy/Mortgage)":
		promptUserForInput(chatID, "Enter category (Rent/Buy/Mortgage):")
	case "Building Age Range":
		promptUserForInput(chatID, "Enter building age range (e.g., 0-20 years):")
	case "Property Type (Apartment/Villa)":
		promptUserForInput(chatID, "Enter property type (Apartment/Villa):")
	case "Floor Range":
		promptUserForInput(chatID, "Enter floor range (e.g., 1-10):")
	case "Storage Availability":
		promptUserForInput(chatID, "Enter storage availability (Yes/No):")
	case "Elevator Availability":
		promptUserForInput(chatID, "Enter elevator availability (Yes/No):")
	case "Advertisement Creation Date Range":
		promptUserForInput(chatID, "Enter advertisement creation date range (e.g., YYYY-MM-DD to YYYY-MM-DD):")
	default:
		sendMessage(int(chatID), "Invalid filter selection.")
		return
	}

	// Temporarily save the selected filter type
	userFilters[userID]["lastFilter"] = selectedFilter
	log.Printf("Saved temporary filter for user %d: %s", userID, selectedFilter)
}

// Function to save user filter input in memory
func saveUserFilterInput(chatId int, userID uint, value string) {
	// Get or initialize the user's FilterItem
	filterItem, exists := userFilterItems[userID]
	if !exists {
		filterItem = &models.FilterItem{}
		userFilterItems[userID] = filterItem
	}

	// Update the relevant field based on filterType
	switch userFilters[uint64(userID)]["lastFilter"] {
	case "Price Range":
		// Assuming the format is "min-max"
		var priceMin, priceMax float64
		if _, err := fmt.Sscanf(value, "%f-%f", &priceMin, &priceMax); err == nil {
			filterItem.PriceMin = priceMin
			filterItem.PriceMax = priceMax
		}
	case "City":
		filterItem.City = value
	case "Neighborhood":
		filterItem.Neighborhood = value
	case "Area Range":
		// Assuming the format is "min-max"
		var areaMin, areaMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &areaMin, &areaMax); err == nil {
			filterItem.AreaMin = areaMin
			filterItem.AreaMax = areaMax
		}
	case "Bedroom Count Range":
		// Assuming the format is "min-max"
		var bedroomsMin, bedroomsMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &bedroomsMin, &bedroomsMax); err == nil {
			filterItem.BedroomsMin = bedroomsMin
			filterItem.BedroomsMax = bedroomsMax
		}
	case "Category (Rent/Buy/Mortgage)":
		filterItem.Category = value
	case "Building Age Range":
		// Assuming the format is "min-max"
		var ageMin, ageMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &ageMin, &ageMax); err == nil {
			filterItem.AgeMin = ageMin
			filterItem.AgeMax = ageMax
		}
	case "Property Type (Apartment/Villa)":
		filterItem.PropertyType = value
	case "Floor Range":
		// Assuming the format is "min-max"
		var floorMin, floorMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &floorMin, &floorMax); err == nil {
			filterItem.FloorMin = floorMin
			filterItem.FloorMax = floorMax
		}
	case "Storage Availability":
		filterItem.HasStorage = (value == "yes")
	case "Elevator Availability":
		filterItem.HasElevator = (value == "yes")
	case "Advertisement Creation Date Range":
		// Assuming the format is "YYYY-MM-DD to YYYY-MM-DD"
		var startDate, endDate time.Time
		dates := strings.Split(value, " to ")
		if len(dates) == 2 {
			startDate, _ = time.Parse("2006-01-02", dates[0])
			endDate, _ = time.Parse("2006-01-02", dates[1])
			filterItem.CreatedDateStart = startDate
			filterItem.CreatedDateEnd = endDate
		}
	}

	sendFilterConfirmationMenu(int64(chatId), strconv.Itoa(int(userFilterItems[userID].ID)))
	// Log the updated filter item
	log.Printf("Updated filter item for user %d: %+v", userID, filterItem)
}

func promptUserForInput(chatID int64, prompt string) {
	sendMessage(int(chatID), prompt)
}

func sendFilterConfirmationMenu(chatID int64, filter string) {
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{
					Text: "Confirm",
					Data: "confirm_filter",
				},
				{
					Text: "Cancel",
					Data: "cancel_filter",
				},
			},
		},
	}

	text := fmt.Sprintf("You selected the filter: *%s*.\nDo you want to confirm or cancel?", filter)
	sendMessageWithInlineKeyboard(int(chatID), text, keyboard)
}

func showFilterMenu(chatID int, userId uint) {
	// Fetch filters from the database
	filters, _ := filterRepository.FindByUserID(userId)
	// Create keyboard buttons for each filter
	var filterButtons [][]KeyboardButton

	// Add the "Create New Filter" button
	filterButtons = append(filterButtons, []KeyboardButton{
		{Text: "Create New Filter"},
	})

	// Proceed safely with filterItems
	if len(filters) == 0 {
		log.Printf("No filters returned for user %d", userId)
	} else {
		for _, filter := range filters {
			filterButtons = append(filterButtons, []KeyboardButton{
				{Text: strconv.Itoa(int(filter.ID))},
			})
		}
	}

	// Define the keyboard layout
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard:        filterButtons,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	// Send the menu
	sendMessageWithKeyboard(int(chatID), "Select a filter or create a new one:", keyboard)
}

func showFilterOptions(chatID int) {
	filterOptions := []string{
		"Price Range",
		"City",
		"Neighborhood",
		"Area Range",
		"Bedroom Count Range",
		"Category (Rent/Buy/Mortgage)",
		"Building Age Range",
		"Property Type (Apartment/Villa)",
		"Floor Range",
		"Storage Availability",
		"Elevator Availability",
		"Advertisement Creation Date Range",
	}

	msg := "Select a filter to apply:"
	sendMessageWithInlineKeyboard(
		int(chatID),
		msg,
		createInlineKeyboardFromOptions(filterOptions),
	)
}

func handleFilterSelection(userID uint, filterID uint) {

	updatedFields := map[string]interface{}{
		"LastFilterItemID": filterID,
	}

	userRepository.UpdateUser(userID, updatedFields)
}

func createFilter(userId uint) {
	// Save the FilterItem
	createdFilterItem, err := filterRepository.Create(*userFilterItems[userId])
	if err != nil {
		fmt.Println("Error saving filter item:", err)
		return
	}

	handleFilterSelection(userId, createdFilterItem.ID)
}

func cancelFilter(userId uint) {
	userFilterItems[userId] = nil
}

func sendFile(chatID int64, content []byte, fileType string) ([]byte, error) {
	var (
		buf    = new(bytes.Buffer)
		writer = multipart.NewWriter(buf)
	)

	chatIdField, err := writer.CreateFormField("chat_id")
	if err != nil {
		return []byte{}, err
	}
	chatIdByteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(chatIdByteArray, uint64(chatID))
	_, err = chatIdField.Write(chatIdByteArray)
	if err != nil {
		return []byte{}, err
	}

	part, err := writer.CreateFormFile("document", "result"+fileType)
	if err != nil {
		return []byte{}, err
	}

	_, err = part.Write(content)
	if err != nil {
		return []byte{}, err
	}

	err = writer.Close()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sendDocument", apiURL), buf)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return cnt, nil
}
