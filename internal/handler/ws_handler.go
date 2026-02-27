// Package handler provides the WebSocket handler for real-time comment broadcasts.
package handler

import (
	"my-portfolio/internal/hub"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

// WSUpgrade is a middleware that ensures the request is a WebSocket upgrade.
// Apply it on the route prefix before the actual WebSocket handler.
func WSUpgrade() fiber.Handler {
	return func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

// WSComments streams real-time comment events over WebSocket.
// Clients receive JSON messages shaped like:
//
//	{"type":"comment","data":{"id":1,"user":"Alice","refresh":true}}
//	{"type":"shutdown","data":{"message":"Server is restarting…"}}
func WSComments(h *hub.Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		ch := h.Register()
		defer h.Unregister(ch)

		for msg := range ch {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	})
}
