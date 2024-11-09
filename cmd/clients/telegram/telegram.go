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
func (c *Client) Updates(offset, limit int)([]Update,error){
	q:=url.Values{}
	q.Add("offset",strconv.Itoa(offset))
	q.Add("limit",strconv.Itoa(limit))

	// ToDo do request
}


func (c *Client)SendMessage(){

}