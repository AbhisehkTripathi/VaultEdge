package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

var HubInstance = NewHub()

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// func HandleWebSocket(c *fiber.Ctx) error {
// 	if websocket.IsWebSocketUpgrade(c) {
// 		return websocket.New(func(conn *websocket.Conn) {
// 			client := &Client{conn: conn, send: make(chan []byte, 256)}
// 			HubInstance.register <- client

// 			go func() {
// 				for {
// 					mt, message, err := conn.ReadMessage()
// 					if err != nil {
// 						break
// 					}
// 					if mt == websocket.TextMessage || mt == websocket.BinaryMessage {
// 						log.Printf("recv: %s", message)
// 					}
// 				}
// 				HubInstance.unregister <- client
// 				conn.Close()
// 			}()

//				// Write pump
//				for msg := range client.send {
//					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
//						break
//					}
//				}
//			})(c)
//		}
//		return fiber.ErrUpgradeRequired
//	}
func HandleWebSocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}
	return websocket.New(func(conn *websocket.Conn) {
		client := &Client{conn: conn, send: make(chan []byte, 256)}
		HubInstance.register <- client

		var closeOnce sync.Once
		cleanup := func() {
			closeOnce.Do(func() {
				HubInstance.unregister <- client
				_ = conn.Close()
			})
		}

		go func() {
			defer cleanup()
			for msg := range client.send {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					break
				}
			}
		}()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if mt == websocket.TextMessage || mt == websocket.BinaryMessage {
				log.Printf("recv: %s", message)
			}
		}
		cleanup()
	})(c)
}

func BroadcastMessage(message []byte) {
	HubInstance.broadcast <- message
}
