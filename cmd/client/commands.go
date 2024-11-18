package client

import (
	"fmt"
	"log"
	"strconv"

	"github.com/MagicalCrawler/RealEstateApp/models"
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
		"Create New Filter":   &CreateFilterCommand{},
		"Location Attachment": &GetLocationAttachmentCommand{},
		"confirm_filter":      &ConfirmFilterCommand{},
		"cancel_filter":         &CancelFilterCommand{},

		"Select Resource Website": &GetResourceWebsite{},
		"Bookmark":                &BookmarkCommand{},
		"Get Website":             &GetWebsiteCommand{},
		"Get Bookmark Id":         &GetBookmarkIDCommand{},
		"Export CSV":              &ExportCSVCommand{},

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

// //////////////////////////////////
type CreateFilterCommand struct{}

func (cmd *CreateFilterCommand) Execute(message *Message, user *models.User) {
	showFilterOptions(message.Chat.ID)
}
func (cmd *CreateFilterCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

// /////////////////////////////////// User Commands
type GetResourceWebsite struct{}

func (cmd *GetResourceWebsite) Execute(message *Message, user *models.User) {

	msg := "Enter \n      s : Sheypoor\n      d : Divar"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}

func (cmd *GetResourceWebsite) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

// /////////////////////////////////// User Commands
type ExportCSVCommand struct{}

func (cmd *ExportCSVCommand) Execute(message *Message, user *models.User) {

	msg := "you enter csv"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}

func (cmd *ExportCSVCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

// ///////////////////////////////////
type GetWebsiteCommand struct{}

func (cmd *GetWebsiteCommand) Execute(message *Message, user *models.User) {

	msg := fmt.Sprintf("Now all your searches will be from the ' %s 'website.", message.Value)
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}

func (cmd *GetWebsiteCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

// ////////////////////////////////////
type StartCommand struct{}

func (cmd *StartCommand) Execute(message *Message, user *models.User) {
	sendMessageWithKeyboard(message.Chat.ID, getWelcomeMessage(message.From.FirstName, user.Role), getKeyboard(user.Role))
}

func (cmd *StartCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER, models.ADMIN, models.SUPER_ADMIN}
}

////// confirm filter

type ConfirmFilterCommand struct{}

func (cmd *ConfirmFilterCommand) Execute(message *Message, user *models.User) {

	msg := fmt.Sprintf("filter confirmed")
	createFilter(user.ID)
	sendMessage(message.Chat.ID, msg)
	return
}

func (cmd *ConfirmFilterCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// ////// cancel filter
type CancelFilterCommand struct{}

func (cmd *CancelFilterCommand) Execute(message *Message, user *models.User) {
	// Here, we would typically remove the user's filter or mark it as canceled
	cancelFilter(user.ID)
	sendMessageWithKeyboard(message.Chat.ID, getWelcomeMessage(message.From.FirstName, user.Role), getKeyboard(user.Role))

	msg := fmt.Sprintf("Your filter has been canceled")
	sendMessage(message.Chat.ID, msg)
	return
}

func (cmd *CancelFilterCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
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

// /////////////////////////////////
type BookmarkCommand struct{}

func (cmd *BookmarkCommand) Execute(message *Message, user *models.User) {
	bookmarks, err := bookmarkRepository.FindAll(uint(message.From.ID))
	var msg string
	if err != nil {
		log.Printf("Error finding bookmarks: %v", err)
		msg = "There was an error fetching your bookmarks. Please try again later."
	} else {
		// msg = "Your bookmarks:\n"
		for _, b := range bookmarks {
			msg += fmt.Sprintf("   ID: %d, Post ID: %+v\n", b.ID, b.Post)
			msg += fmt.Sprint("-------------------------\n")
		}
	}

	msg += "\nTo create new bookmark:\n\tsend me Post ID with pattern 'B=<number>'"
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *BookmarkCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// /////////////////////////////////
type GetBookmarkIDCommand struct{}

func (cmd *GetBookmarkIDCommand) Execute(message *Message, user *models.User) {
	var msg string
	id, err := strconv.ParseInt(message.Value[2:], 10, 64)
	if err != nil {
		msg = "Invalid ID format. Please use 'B=<number>'."
	}
	post, err := postRepository.FindByID(uint(id))
	if err != nil {
		log.Printf("Error finding post: %v", err)
		msg = "There was an error fetching post. Please try again later."
	} else {
		user, err := userRepository.Find(uint(message.From.ID))
		if err != nil {
			log.Printf("Error finding user: %v", err)
			msg = "There was an error fetching user. Please try again later."
		} else {
			err = bookmarkRepository.Save(post, user)
			if err != nil {
				log.Printf("Error saving bookmark: %v", err)
				msg = "There was an error bookmarking this post. Please try again later."
			} else {
				msg = "Done!"
			}
		}
	}
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
}
func (cmd *GetBookmarkIDCommand) AllowedRoles() []models.Role {
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
	// filterOptions := []string{
	// 	"Price Range",
	// 	"City",
	// 	"Neighborhood",
	// 	"Area Range",
	// 	"Bedroom Count Range",
	// 	"Category (Rent/Buy/Mortgage)",
	// 	"Building Age Range",
	// 	"Property Type (Apartment/Villa)",
	// 	"Floor Range",
	// 	"Storage Availability",
	// 	"Elevator Availability",
	// 	"Advertisement Creation Date Range",
	// }

	// msg := "Select a filter to apply:"
	showFilterMenu(message.Chat.ID, user.ID)
	// sendMessageWithInlineKeyboard(message.Chat.ID, msg, createInlineKeyboardFromOptions(filterOptions))
}

func (cmd *FilterCommand) AllowedRoles() []models.Role {
	return []models.Role{models.USER}
}

// ////////////////////////////////
type PopularsCommand struct{}

func (cmd *PopularsCommand) Execute(message *Message, user *models.User) {
	msg := "    All Popular Advertisements:\n\n"
	ads, err := postRepository.GetMostVisitedPost()
	if err != nil {
		msg = "Error fetching posts: "
		log.Fatal(msg + err.Error())
	} else {
		if len(ads) == 0 {
			msg = "Nothing found"
		} else {
			for _, a := range ads {
				msg += fmt.Sprintf("%s  \n", a.Title) //, a.Post.WatchedNum)
				msg += fmt.Sprint("-------------------------\n")
			}
		}
	}
	sendMessageWithKeyboard(message.Chat.ID, msg, getKeyboard(user.Role))
	return
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
	msg := fmt.Sprintf("User with id :%s changed to Premium client.", message.Value[3:])
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

	msg := "You entered Monitor\nCrawls\n\n"
	/////////////
	crawlHistories := postRepository.GetAllCrawlHistory()
	for _, crawl := range crawlHistories {
		msg += fmt.Sprintf("\nID: %v, CPU: %v, Memory: %v\n", crawl.ID, crawl.CpuUsage, crawl.MemoryUsage)

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
	msg := "    Advertisements:\n\n"
	ads, err := postRepository.GetAllPosts()
	if err != nil {
		msg = "Error fetching posts: "
		log.Fatal(msg + err.Error())
	} else {
		if len(ads) == 0 {
			msg = "Nothing found"
		} else {
			for _, a := range ads {
				msg += fmt.Sprintf("%s  \n", a.Title) //, a.Post.WatchedNum)
				msg += fmt.Sprint("-------------------------\n")
			}
		}
	}
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
