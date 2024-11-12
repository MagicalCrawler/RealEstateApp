package main

import (
	telegram_client "github.com/MagicalCrawler/RealEstateApp/cmd/client"
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {
	utils.LoadEnvFile()
	db.NewConnection()

	telegram_client.Run()

}
