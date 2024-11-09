package telegram

const BASE_PATH = "https://api.telegram.org/bot"

type Client struct{
	host     string
	basePath string
	client http.Client
}

