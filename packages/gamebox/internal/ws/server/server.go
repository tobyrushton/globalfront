package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/tobyrushton/globalfront/packages/gamebox/internal/client"
)

type WsServer struct {
	clientsMu sync.Mutex
	clients   map[*client.Client]struct{}
}

func NewServer() *WsServer {
	return &WsServer{
		clients: make(map[*client.Client]struct{}),
	}
}

func (s *WsServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer c.CloseNow()

	cl := client.New(c)

	s.clientsMu.Lock()
	s.clients[cl] = struct{}{}
	s.clientsMu.Unlock()

	// read loop
	go func() {
		for {
			_, msg, err := c.Read(context.Background())
			if err != nil {
				break
			}
			// do something with msg
			fmt.Println(msg)
		}
	}()

	s.clientsMu.Lock()
	delete(s.clients, cl)
	s.clientsMu.Unlock()
}
