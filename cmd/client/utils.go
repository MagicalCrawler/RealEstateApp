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
				}

				if update.Callback != nil {
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
	user, err := userRepository.FindByTelegramID(userID)

	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return
	}

	selectedFilter := callbackQuery.Data
	chatID := int64(callbackQuery.Message.Chat.ID)

	// Initialize the user's filter map if it doesn't exist
	if _, exists := userFilters[uint64(user.ID)]; !exists {
		userFilters[uint64(user.ID)] = make(map[string]string)
	}

	if strings.HasPrefix(callbackQuery.Data, "post_") {
		// Extract the post ID
		postIDStr := strings.TrimPrefix(callbackQuery.Data, "post_")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			sendMessage(int(chatID), "Invalid post selection.")
			return
		}

		sendMessage(postID, "post id")
		return
	}

	if strings.HasPrefix(callbackQuery.Data, "filter_") {
		// Extract the filter ID
		filterIDStr := strings.TrimPrefix(callbackQuery.Data, "filter_")
		filterID, err := strconv.Atoi(filterIDStr)
		if err != nil {
			sendMessage(int(chatID), "Invalid filter selection.")
			return
		}

		// Use the filterID as needed
		sendMessageWithKeyboard(int(chatID), fmt.Sprintf("Selected filter ID: %d", filterID), getKeyboard(user.Role))

		handleFilterSelection(user.ID, uint(filterID))
		return
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
	userFilters[uint64(user.ID)]["lastFilter"] = selectedFilter
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

	// Debug: Log the filter type and value
	filterType := userFilters[uint64(userID)]["lastFilter"]
	log.Printf("Filter Type: %s, Value: %s", filterType, value)

	// Update the relevant field based on filterType
	switch filterType {
	case "Price Range":
		var priceMin, priceMax float64
		if _, err := fmt.Sscanf(value, "%f-%f", &priceMin, &priceMax); err == nil {
			filterItem.PriceMin = priceMin
			filterItem.PriceMax = priceMax
		} else {
			log.Printf("Error parsing Price Range: %v", err)
		}
	case "City":
		filterItem.City = value
		log.Printf("City updated to: %s", filterItem.City) // Debug: Log the update
	case "Neighborhood":
		filterItem.Neighborhood = value
	case "Area Range":
		var areaMin, areaMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &areaMin, &areaMax); err == nil {
			filterItem.AreaMin = areaMin
			filterItem.AreaMax = areaMax
		} else {
			log.Printf("Error parsing Area Range: %v", err)
		}
	case "Bedroom Count Range":
		var bedroomsMin, bedroomsMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &bedroomsMin, &bedroomsMax); err == nil {
			filterItem.BedroomsMin = bedroomsMin
			filterItem.BedroomsMax = bedroomsMax
		} else {
			log.Printf("Error parsing Bedroom Count Range: %v", err)
		}
	case "Category (Rent/Buy/Mortgage)":
		filterItem.Category = value
	case "Building Age Range":
		var ageMin, ageMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &ageMin, &ageMax); err == nil {
			filterItem.AgeMin = ageMin
			filterItem.AgeMax = ageMax
		} else {
			log.Printf("Error parsing Building Age Range: %v", err)
		}
	case "Property Type (Apartment/Villa)":
		filterItem.PropertyType = value
	case "Floor Range":
		var floorMin, floorMax int
		if _, err := fmt.Sscanf(value, "%d-%d", &floorMin, &floorMax); err == nil {
			filterItem.FloorMin = floorMin
			filterItem.FloorMax = floorMax
		} else {
			log.Printf("Error parsing Floor Range: %v", err)
		}
	case "Storage Availability":
		filterItem.HasStorage = (value == "yes")
	case "Elevator Availability":
		filterItem.HasElevator = (value == "yes")
	case "Advertisement Creation Date Range":
		var startDate, endDate time.Time
		dates := strings.Split(value, " to ")
		if len(dates) == 2 {
			startDate, _ = time.Parse("2006-01-02", dates[0])
			endDate, _ = time.Parse("2006-01-02", dates[1])
			filterItem.CreatedDateStart = startDate
			filterItem.CreatedDateEnd = endDate
		}
	}

	// Send confirmation menu
	sendFilterConfirmationMenu(int64(chatId))

	// Log the updated filter item
	log.Printf("Updated filter item for user %d: %+v", userID, filterItem)
}

func promptUserForInput(chatID int64, prompt string) {
	sendMessage(int(chatID), prompt)
}

func sendFilterConfirmationMenu(chatID int64) {
	keyboard := ReplyKeyboardMarkupWithLocation{
		Keyboard: [][]KeyboardButton{
			{
				KeyboardButton{
					Text: "SaveFilter",
				},
				KeyboardButton{
					Text: "CancelFilter",
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	sendMessageWithKeyboard(int(chatID), "continue add filter", keyboard)
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
				{Text: "filter_" + strconv.Itoa(int(filter.ID))},
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
	sendFilterConfirmationMenu(int64(chatID))
}

func handleFilterSelection(userID uint, filterID uint) {

	updatedFields := map[string]interface{}{
		"LastFilterItemID": filterID,
	}

	userRepository.UpdateUser(userID, updatedFields)
}

func createFilter(userId uint) {
	// Save the FilterItem
	userFilterItems[userId].UserID = userId
	createdFilterItem, err := filterRepository.Create(*userFilterItems[userId])

	if err != nil {
		fmt.Println("Error saving filter item:", err)
		return
	}

	handleFilterSelection(userId, createdFilterItem.ID)
	cancelFilter(userId)
}

func cancelFilter(userId uint) {
	userFilterItems[userId] = nil
	userFilters[uint64(userId)] = nil
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

func searchLastFilter(chatID int, userID int) {
	lastFilterItem, err := userRepository.GetLastFilterItem(uint(userID))
	if err != nil {
		log.Printf("Error retrieving last filter item: %v", err)
		return
	}
	if lastFilterItem == nil {
		log.Printf("Please add some filter :)")
		return
	}

	posts, err := filterRepository.SearchPostHistory(*lastFilterItem)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		sendMessage(chatID, "An error occurred while fetching posts.")
		return
	}

	if len(posts) == 0 {
		// Inform the user if no posts match the filter
		sendMessage(chatID, "No posts found matching your filter.")
		return
	}

	// Generate the inline keyboard
	keyboard := createInlineKeyboardFromPosts(posts)

	// Prepare the message text
	text := "Here are the posts matching your filter:\nSelect one to view details."

	// Send the message with the inline keyboard
	sendMessageWithInlineKeyboard(chatID, text, keyboard)
}

// Function to create an inline keyboard from a list of posts
func createInlineKeyboardFromPosts(posts []models.PostHistory) InlineKeyboardMarkup {
	buttons := make([][]InlineKeyboardButton, 0)

	for _, post := range posts {
		text := fmt.Sprintf(
			"🏡 *%s*\n\nPrice: %d\nCity: %s\nNeighborhood: %s\nArea: %d m²\nBedrooms: %d\n\n[View Post](%s)",
			post.Title, post.Price, post.City, post.Neighborhood, post.Area, post.BedroomNum, post.PostURL,
		)
		// Each row will have one button with the post's title or summary
		row := []InlineKeyboardButton{
			{
				Text: text,
				// Text: fmt.Sprintf("%s - %d", post.Title, post.Price), // Display title and price
				Data: fmt.Sprintf("post_%d", post.ID), // Unique callback data
			},
		}
		buttons = append(buttons, row)
	}

	return InlineKeyboardMarkup{InlineKeyboard: buttons}
}
