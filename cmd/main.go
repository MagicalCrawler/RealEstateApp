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

	// Initialize crawler service jobs
	crawlerService := services.NewCrawlerService()
	crawlerService.Start()

	// Run the Telegram bot
	client.Run(dbConnection)

	select {}
}
