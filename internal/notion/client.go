package notion

import (
	"github.com/jomei/notionapi"
)

type Client struct {
	api *notionapi.Client
}

func NewClient(token string) *Client {
	return &Client{
		api: notionapi.NewClient(notionapi.Token(token)),
	}
}

func (c *Client) API() *notionapi.Client {
	return c.api
}
