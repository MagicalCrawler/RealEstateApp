package main

import (
	"github.com/MagicalCrawler/RealEstateApp/cmd/client"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/services"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {

	// Load environment variables
	utils.LoadEnvFile()

	// Initialize DB connection
	dbConnection := db.NewConnection()

	// Run the Telegram bot
	client.Run(dbConnection)

	go services.StartCrawlers()

	select {}
}
