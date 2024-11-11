package main

import (
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {
	utils.LoadEnvFile()
	db.NewConnection()
}
