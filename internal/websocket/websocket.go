package websocket

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return websocket.New(func(c *websocket.Conn) {
			// c.Locals is added to the *websocket.Conn
			log.Println(c.Locals("allowed"))  // true
			log.Println(c.Params("id"))       // 123
			log.Println(c.Query("v"))         // 1.0
			log.Println(c.Cookies("session")) // ""

			// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
			var (
				mt  int
				msg []byte
				err error
			)
			for {
				if mt, msg, err = c.ReadMessage(); err != nil {
					log.Println("read:", err)
					break
				}
				log.Printf("recv: %s", msg)

				// Echo the message back
				if err = c.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

		})(c)
	}
	// Returns status 426 Upgrade Required
	return fiber.ErrUpgradeRequired
}

// BroadcastMessage broadcasts a message to all connected clients
func BroadcastMessage(message []byte) {
	// In a real application, you would maintain a list of active connections
	// and broadcast to all of them
	log.Printf("Broadcasting message: %s", string(message))
}
