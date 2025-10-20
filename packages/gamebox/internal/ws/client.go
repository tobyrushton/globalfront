package ws

import (
	"github.com/gorilla/websocket"
	pb "github.com/tobyrushton/globalfront/pb/messages/v1"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn        *websocket.Conn
	sendChannel chan *pb.WebsocketMessage

	playerId string
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		conn:        c,
		sendChannel: make(chan *pb.WebsocketMessage),
	}
}

func (c *Client) Send(message *pb.WebsocketMessage) error {
	m, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, m)
}

func (c *Client) GetSendChannel() chan<- *pb.WebsocketMessage {
	return c.sendChannel
}

func (c *Client) PlayerId() string {
	return c.playerId
}

func (c *Client) SetPlayerId(id string) {
	c.playerId = id
}
