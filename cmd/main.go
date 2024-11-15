package main

import (
	"log/slog"

	"github.com/MagicalCrawler/RealEstateApp/cmd/client"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/services"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {

	slog.Info("Load environment variables")
	utils.LoadEnvFile()

	logger := utils.MainLogger()

	logger.Debug("Initialize DB connection")
	dbConnection := db.NewConnection()

  logger.Debug("Initialize crawler service jobs")
	crawlerService := services.NewCrawlerService()
	crawlerService.Start()

  logger.Debug("Run the Telegram bot")
	client.Run(dbConnection)

	select {}
}
