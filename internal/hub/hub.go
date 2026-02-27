// Package hub implements a WebSocket broadcast hub.
// It is safe for concurrent use. Multiple goroutines may register/unregister
// clients and broadcast events simultaneously.
package hub

import (
	"encoding/json"
	"sync"
)

// EventType classifies events sent to connected clients.
type EventType string

const (
	// EventComment fires when a new approved comment is available.
	EventComment EventType = "comment"
	// EventShutdown fires when the server is about to shut down.
	EventShutdown EventType = "shutdown"
)

// Event is the payload broadcast to every connected WebSocket client.
type Event struct {
	Type EventType `json:"type"`
	Data any       `json:"data"`
}

// client is a buffered channel for sending JSON-encoded event bytes.
type client chan []byte

// Hub manages all active WebSocket client channels.
type Hub struct {
	mu      sync.RWMutex
	clients map[client]struct{}
}

// New creates and returns an initialized *Hub.
func New() *Hub {
	return &Hub{clients: make(map[client]struct{})}
}

// Register adds a new client channel to the hub and returns it.
// The caller is responsible for calling Unregister when done.
func (h *Hub) Register() client {
	ch := make(client, 8) // buffered to avoid blocking on slow clients
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

// Unregister removes a client channel from the hub and closes it.
func (h *Hub) Unregister(ch client) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

// Broadcast sends a JSON-encoded event to every connected client.
// Clients whose channels are full are skipped (non-blocking).
func (h *Hub) Broadcast(ev Event) {
	data, err := json.Marshal(ev)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- data:
		default: // skip lagging client
		}
	}
}

// ClientCount returns the number of currently registered clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

