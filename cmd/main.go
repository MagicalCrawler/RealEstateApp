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
	userRepository := db.CreateNewUserRepository(dbConnection)
	postRepository := db.NewPostRepository(dbConnection)
	crawlerService := services.NewCrawlerService(&postRepository)
	crawlerService.Start()

	logger.Debug("Run the Telegram bot")
	client.Run(userRepository, postRepository)
}
