package client

import "github.com/coder/websocket"

type Client struct {
	conn *websocket.Conn
}

func New(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
	}
}
