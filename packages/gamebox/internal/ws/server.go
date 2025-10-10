package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	pb "github.com/tobyrushton/globalfront/pb/messages/v1"
	"google.golang.org/protobuf/proto"
)

type WsServer struct {
	clientsMu sync.Mutex
	clients   map[string]*Client

	msgChan chan *pb.WebsocketMessage
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer(msgChan chan *pb.WebsocketMessage, playerIds []string) *WsServer {
	clients := make(map[string]*Client)
	for _, id := range playerIds {
		clients[id] = nil
	}
	return &WsServer{
		clients: clients,
		msgChan: msgChan,
	}
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer c.Close()

	cl := NewClient(c)
	done := make(chan struct{})

	// read loop
	go func() {
		msg, err := s.readMsg(c)
		if err != nil {
			close(done)
			return
		}
		switch v := msg.Payload.(type) {
		case *pb.WebsocketMessage_JoinGame:
			err := s.addClient(v.JoinGame.PlayerId, cl)
			// TODO: Send error message back to client
			if err != nil {
				fmt.Println("Error adding client:", err)
				close(done)
				return
			}
			fmt.Println("Player joined:", v.JoinGame.PlayerId)
			s.msgChan <- msg
		default:
			fmt.Println("Expected JoinGame message, got:", v)
			return
		}

		for {
			msg, err := s.readMsg(c)
			if err != nil {
				close(done)
				return
			}
			s.msgChan <- msg
		}
	}()

	for {
		select {
		case msg := <-cl.sendChannel:
			err := cl.Send(msg)
			if err != nil {
				return
			}
		case <-done:
			s.clientsMu.Lock()
			s.clients[cl.PlayerId()] = nil
			s.clientsMu.Unlock()
		}
	}
}

func (s *WsServer) Broadcast(message *pb.WebsocketMessage) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, cl := range s.clients {
		if cl != nil {
			cl.GetSendChannel() <- message
		}
	}
}

func (s *WsServer) readMsg(c *websocket.Conn) (*pb.WebsocketMessage, error) {
	_, data, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	var msg pb.WebsocketMessage
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *WsServer) addClient(playerId string, cl *Client) error {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	if client, exists := s.clients[playerId]; !exists {
		return fmt.Errorf("player ID %s not recognized", playerId)
	} else if client != nil {
		return fmt.Errorf("player ID %s already has a connected client", playerId)
	} else {
		s.clients[playerId] = cl
	}
	return nil
}
