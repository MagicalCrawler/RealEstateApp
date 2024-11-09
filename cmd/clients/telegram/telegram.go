package telegram

const BASE_PATH = "https://api.telegram.org/bot"

type Client struct{
	host     string
	basePath string
	client http.Client
}

func New(host, token string) Client{
	return Client{
		host:host,
		basePath:BASE_PATH + token
		basePath:http.Client{}
	}
}
