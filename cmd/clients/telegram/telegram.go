package telegram

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/MagicalCrawler/RealEstateApp/cmd/err"

	
)

const (
	BASE_PATH ="https://api.telegram.org/bot"
	getUpdatesMethod="getUpdates"
	sendMessageMethod="sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}


func New(host, token string) *Client{
	return &Client{
		host:    host,
		basePath:BASE_PATH + token,
		client:  http.Client{},
	}
}
func (c *Client) Updates(ctx context.Context, offset int, limit int)([]Update,error){
	q:=url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data,er:=c.doRequest(ctx,getUpdatesMethod,q)
	if er != nil{
		return nil, er
	}
	var res UpdatesResponse

	if er:=json.Unmarshal(data,&res);er!=nil{
		return nil,er
	}
	return res.Result,nil
}

func (c *Client)SendMessage(ctx context.Context,chatID int, text string)error{
	q:=url.Values{}
	q.Add("chat_id",strconv.Itoa(chatID))
	q.Add("text",text)

	_,er:=c.doRequest(ctx,sendMessageMethod,q)
	if er != nil{
		return err.Wrap("can't send message",er)
	}
	return nil
}

func (c *Client) doRequest(ctx context.Context,method string, query url.Values)([]byte,error){
	// defer func(){er=err.WrapIfErr("can't do request",er)}()
	u:=url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, er := http.NewRequestWithContext(ctx,http.MethodGet, u.String(), nil)
	if er != nil{
		return nil, er
	}
	req.URL.RawQuery=query.Encode()
	resp,er:=c.client.Do(req)
	if er != nil{
		return nil, er
	}
	defer func() {_=resp.Body.Close()}()
	
	body,er:=io.ReadAll(resp.Body)
	if er != nil{
		return nil, er
	}
	return body,nil
}

