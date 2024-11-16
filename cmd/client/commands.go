package client

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/MagicalCrawler/RealEstateApp/models"
	"gorm.io/gorm"
)

func initializeCommands() {
	CommandRegistry = map[string]Command{
		"/start": &StartCommand{},
		//user commands
		"Help":          &HelpCommand{},
		"Send Location": &SendLocationCommand{},
		"Search":        &SearchCommand{},
		"Populars":      &PopularsCommand{},
		"Get Redius":    &GetRediusCommand{},

		"Setting":             &SettingCommand{},
		"Filter":              &FilterCommand{},
		"Location Attachment": &GetLocationAttachmentCommand{},
		//admin commands
		"Premium":           &PremiumCommand{},
		"Errors":            &ErrorsCommand{},
		"Clients":           &ClientCommand{},  //admin and super-admin
		"Filters":           &FiltersCommand{}, //
		"Change To Premium": &GetPremiumIdCommand{},
		//super-admin commads
		"Admins":          &AdminCommand{},
		"Get Admin Id":    &GetAdminIdCommand{},
		"Create Admin":    &CreateAdminCommand{},
		"Monitor":         &MonitorCommand{},
		"Advertisements":  &AdvertisementsCommand{},
		"Crawler Setting": &CrawlerSettingCommand{},
	}
}

// /////////////////////////////////// User Commands
type StartCommand struct{}

func (cmd *StartCommand) Execute(message *Message, user *models.User) {
	sendMessageWithKeyboard(message.Chat.ID, getWelcomeMessage(message.From.FirstName, user.Role), getKeyboard(user.Role))
}

func (cmd *StartCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

//////////////////////////////////////

type GetRediusCommand struct{}

func (cmd *GetRediusCommand) Execute(message *Message, user *models.User) {

	msg := fmt.Sprintf("You entered radius: %s", message.Value[7:])
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *GetRediusCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// //////////////////////////////////
type GetLocationAttachmentCommand struct{}

func (cmd *GetLocationAttachmentCommand) Execute(message *Message, user *models.User) {
	msg := fmt.Sprintf("Your selected location is with latitude: %f, and longitude: %f👌\n\nNow send me your desired radius with pattern👉 \"redius=<number>\"", message.Location.Latitude, message.Location.Longitude)
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *GetLocationAttachmentCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// ////////////////////////////////
type SettingCommand struct{}

func (cmd *SettingCommand) Execute(message *Message, user *models.User) {
	msg := "You entered setting"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *SettingCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// //////////////////////////////////

type HelpCommand struct{}

func (cmd *HelpCommand) Execute(message *Message, user *models.User) {
	msg := `Real Estate Finder Bot!
                    /search to find properties based on filters like price, location, and type.
                    /help for more information.`
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}

func (cmd *HelpCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

// /////////////////////////////////
type SendLocationCommand struct{}

func (cmd *SendLocationCommand) Execute(message *Message, user *models.User) {
	msg := "You can send me location📍 with your telegram attachment 👇"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *SendLocationCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// ////////////////////////////////////
type SearchCommand struct{}

func (cmd *SearchCommand) Execute(message *Message, user *models.User) {
	msg := fmt.Sprintf("You entered search ")
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *SearchCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// ////////////////////////////////////
type FilterCommand struct{}

func (cmd *FilterCommand) Execute(message *Message, user *models.User) {
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
	sendMessageWithInlineKeyboard(message.Chat.ID, msg, createInlineKeyboardFromOptions(filterOptions))
}

func (cmd *FilterCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
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

func handleCallbackQuery(callbackQuery *CallbackQuery) {
	userID := uint64(callbackQuery.From.ID)
	_, err := userRepository.FindByTelegramID(userID)

	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return
	}

	selectedFilter := callbackQuery.Data
	chatID := int64(callbackQuery.Message.Chat.ID)
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
		sendMessage(callbackQuery.Message.Chat.ID, "Invalid filter selection.")
	}
}

func saveUserFilterInput(db *gorm.DB, userID uint, filterType, value string) {
	var filterItem models.FilterItem

	// Find existing FilterItem for the user (if any)
	if err := db.Where("user_id = ?", userID).First(&filterItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Initialize a new FilterItem if none exists
			filterItem = models.FilterItem{}
		} else {
			log.Printf("Error retrieving filter item for user %d: %v", userID, err)
			return
		}
	}

	// Update the relevant field based on filterType
	switch filterType {
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

	// Save or update the FilterItem in the database
	if filterItem.ID == 0 {
		// Create a new record
		if _, err := models.CreateFilterItem(db, filterItem); err != nil {
			log.Printf("Error creating filter item for user %d: %v", userID, err)
		}
	} else {
		// Update the existing record
		if err := db.Save(&filterItem).Error; err != nil {
			log.Printf("Error updating filter item for user %d: %v", userID, err)
		}
	}
}

func promptUserForInput(chatID int64, prompt string) {
	sendMessage(int(chatID), prompt)
}

// ////////////////////////////////
type PopularsCommand struct{}

func (cmd *PopularsCommand) Execute(message *Message, user *models.User) {
	msg := "You entered populars"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *PopularsCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

///////////////////////////////////////////////// Admin Commands

type ErrorsCommand struct{}

func (cmd *ErrorsCommand) Execute(message *Message, user *models.User) {
	msg := "You entered errors"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *ErrorsCommand) AllowedRoles() []models.Role {
	return []models.Role{models.ADMIN, models.SUPER_ADMIN}
}

// ///////////////////////////////
type ClientCommand struct{}

func (cmd *ClientCommand) Execute(message *Message, user *models.User) {
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
}
func (cmd *ClientCommand) AllowedRoles() []models.Role {
	return []models.Role{models.ADMIN, models.SUPER_ADMIN}
}

// ////////////////////////////////
type FiltersCommand struct{}

func (cmd *FiltersCommand) Execute(message *Message, user *models.User) {
	msg := "You entered filters"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *FiltersCommand) AllowedRoles() []models.Role {
	return []models.Role{models.ADMIN, models.SUPER_ADMIN}
}

// ///////////////////////////////
type PremiumCommand struct{}

func (cmd *PremiumCommand) Execute(message *Message, user *models.User) {
	msg := "Send me user id with pattern👉 \"Id=<number>\""
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *PremiumCommand) AllowedRoles() []models.Role {
	return []models.Role{models.ADMIN, models.SUPER_ADMIN}
}

// ////////////////////////////////
type GetPremiumIdCommand struct{}

func (cmd *GetPremiumIdCommand) Execute(message *Message, user *models.User) {
	msg := fmt.Sprintf("User with id :%s changed to Premium client.", message.Title)
	id, err := strconv.ParseInt(message.Value[3:], 10, 64)
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
}
func (cmd *GetPremiumIdCommand) AllowedRoles() []models.Role {
	return []models.Role{models.ADMIN, models.SUPER_ADMIN}
}

//////////////////////////////////////////////////// Super-admin commands

type AdminCommand struct{}

func (cmd *AdminCommand) Execute(message *Message, user *models.User) {
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
		msg += "\nEnter 'c' to Create Admin"
	}
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *AdminCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}

// //////////////////////////////////
type MonitorCommand struct{}

func (cmd *MonitorCommand) Execute(message *Message, user *models.User) {
	msg := "You entered Monitor\nCrawls"
	/////////////
	crawlHistories, _ := postRepository.GetCrawlHistory()
	for _, ch := range crawlHistories {
		msg += fmt.Sprintf("\nID: %d, CPU: %v, Memory: %v\n", ch.ID, ch.CpuUsage, ch.MemoryUsage)
	}
	/////////////
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *MonitorCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}

// //////////////////////////////////
type AdvertisementsCommand struct{}

func (cmd *AdvertisementsCommand) Execute(message *Message, user *models.User) {
	msg := "You entered Advertisements"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *AdvertisementsCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}

// //////////////////////////////////
type CrawlerSettingCommand struct{}

func (cmd *CrawlerSettingCommand) Execute(message *Message, user *models.User) {
	msg := "You entered Crawler Setting"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *CrawlerSettingCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}

// //////////////////////////////////
type CreateAdminCommand struct{}

func (cmd *CreateAdminCommand) Execute(message *Message, user *models.User) {
	msg := fmt.Sprintf("User with id :%s changed to Admin.", message.Value[6:])
	id, err := strconv.ParseInt(message.Value[6:], 10, 64)
	if err != nil {
		msg = "Invalid ID format. Please use 'admin=<number>'."
	} else {
		_, err := userRepository.UpdateUserRole(uint(id), models.ADMIN)
		if err != nil {
			msg = fmt.Sprintf("Error updating user role: %v", err)
		}
		// msg += "\nEnter 'c' to Create Admin"
	}
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *CreateAdminCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}

// ///////////////////////////////
type GetAdminIdCommand struct{}

func (cmd *GetAdminIdCommand) Execute(message *Message, user *models.User) {
	msg := "Send me user's ID with pattern \"admin=<number>\""
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
}
func (cmd *GetAdminIdCommand) AllowedRoles() []models.Role {
	return []models.Role{models.SUPER_ADMIN}
}
