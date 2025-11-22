package httpc

import "net/http"

type Config struct{}

type Client struct {
	c *http.Client
}

func New() *Client {
	return &Client{
		c: &http.Client{},
	}
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.c.Do(req)
	return
}
