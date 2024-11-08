package main

import (
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {
	utils.LoadEnvFile()
	db.NewConnection()

	token = utils.GetConfig("TELEGRAM_TOKEN")
	if token==""{
		panic("Telegram token not found in .env file")
	}
	// tgClient = telegram.New(token)
	// fetcher = fetcher.New(tgClient)
	// processor = processor.New(tgClient)
	//consumer.Start(fetcher,processor)

}

