package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type SocketBroadcaster struct {
	mu    sync.RWMutex
	conns []*websocket.Conn
}

func NewSocketBroadcaster() *SocketBroadcaster {
	return &SocketBroadcaster{
		conns: make([]*websocket.Conn, 0),
	}
}

func (b *SocketBroadcaster) Add(conn *websocket.Conn) {
	b.mu.Lock()
	b.conns = append(b.conns, conn)
	b.mu.Unlock()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				b.Remove(conn)
				return
			}
		}
	}()
}

func (b *SocketBroadcaster) Remove(conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, c := range b.conns {
		if c == conn {
			b.conns = append(b.conns[:i], b.conns[i+1:]...)
			break
		}
	}
}

func (b *SocketBroadcaster) Broadcast(message interface{}) error {
	b.mu.RLock()
	conns := make([]*websocket.Conn, len(b.conns))
	copy(conns, b.conns)
	b.mu.RUnlock()

	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			b.Remove(conn)
		}
	}

	return nil
}
