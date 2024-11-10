package main

import (
	"log"
	// "github.com/MagicalCrawler/RealEstateApp/db"
	"github.com/MagicalCrawler/RealEstateApp/utils"
	tgEvent "github.com/MagicalCrawler/RealEstateApp/cmd/events/telegram_event"
	tgClient "github.com/MagicalCrawler/RealEstateApp/cmd/clients/telegram"
	event_consumer "github.com/MagicalCrawler/RealEstateApp/cmd/consumer"
)

const (
	tgBotHost   = "api.telegram.org"
	batchSize   = 100
)

func main() {
	utils.LoadEnvFile()
	// db.NewConnection()

	token := utils.GetConfig("TELEGRAM_TOKEN")
	log.Println("Bot Token:", token)
	eventsProcessor := tgEvent.New(
		tgClient.New(tgBotHost, token),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

