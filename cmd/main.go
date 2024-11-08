package main

import (
	"github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
)

func main() {
	utils.LoadEnvFile()
	db.NewConnection()

	// token = GetToken()
	// tgClient = telegram.New(token)
	// fetcher = fetcher.New(tgClient)
	// processor = processor.New(tgClient)
	//consumer.Start(fetcher,processor)

}
