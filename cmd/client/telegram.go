package client

import (
	"github.com/MagicalCrawler/RealEstateApp/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var botToken = utils.GetConfig("TELEGRAM_TOKEN")
const (
	apiURL   = "https://api.telegram.org/bot" + botToken

	msgHelp = `Real Estate Finder Bot!
				/search to find properties based on filters like price, location, and type.
				/notify to get alerts for new listings matching your preferences.
				/help for more information.`
	batchSize = 100
	timeout = 10
)

log.Print("service started")
func Run(){
	
}