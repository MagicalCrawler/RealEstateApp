package main

import (
	"github.com/MagicalCrawler/RealEstateApp/services"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	//utils.LoadEnvFile()
	//db.NewConnection()
	//
	//token := utils.GetConfig("TELEGRAM_TOKEN")
	//log.Println("Bot Token:", token)
	//eventsProcessor := tgEvent.New(
	//	tgClient.New(tgBotHost, token),
	//)
	//
	//log.Print("service started")
	//
	//consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	//
	//if err := consumer.Start(); err != nil {
	//	log.Fatal("service is stopped", err)
	//}

	go services.StartCrawlers()

	select {}
}
