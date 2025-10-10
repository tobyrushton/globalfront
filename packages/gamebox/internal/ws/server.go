package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	pb "github.com/tobyrushton/globalfront/pb/messages/v1"
	"google.golang.org/protobuf/proto"
)

type WsServer struct {
	clientsMu sync.Mutex
	clients   map[*Client]struct{}

	msgChan chan *pb.WebsocketMessage
}

func NewServer(msgChan chan *pb.WebsocketMessage) *WsServer {
	return &WsServer{
		clients: make(map[*Client]struct{}),
		msgChan: msgChan,
	}
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:3000"},
	})

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer c.CloseNow()

	cl := NewClient(c)

	s.clientsMu.Lock()
	s.clients[cl] = struct{}{}
	s.clientsMu.Unlock()

	// read loop
	go func() {
		for {
			_, data, err := c.Read(context.Background())
			if err != nil {
				break
			}
			var msg pb.WebsocketMessage
			if err := proto.Unmarshal(data, &msg); err != nil {
				fmt.Println("Failed to unmarshal message:", err)
			} else {
				s.msgChan <- &msg
			}
		}
	}()

	ctx := c.CloseRead(context.Background())
	for {
		select {
		case msg := <-cl.sendChannel:
			err := cl.Send(msg)
			if err != nil {
				return
			}
		case <-ctx.Done():
			s.clientsMu.Lock()
			delete(s.clients, cl)
			s.clientsMu.Unlock()
		}
	}
}

func (s *WsServer) Broadcast(message *pb.WebsocketMessage) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for cl := range s.clients {
		cl.GetSendChannel() <- message
	}
}
