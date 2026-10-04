package handlers

import "sync"

// clientHub tracks all connected SSE clients and broadcasts frames to them.
type clientHub struct {
	mu      sync.Mutex
	clients map[uint64]chan []byte
	nextID  uint64
	count   int
}

// Register adds a client and returns its ID + frame channel.
func (h *clientHub) Register() (uint64, chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	id := h.nextID
	h.nextID++
	ch := make(chan []byte, 2)
	h.clients[id] = ch
	h.count++
	return id, ch
}

// Unregister removes a client and closes its channel.
func (h *clientHub) Unregister(id uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.clients[id]; ok {
		close(ch)
		delete(h.clients, id)
		h.count--
	}
}

// PlayerCount returns the number of active players.
// Safe to call while gameMu is held (lock order: gameMu → hubMu).
func (h *clientHub) PlayerCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.count
}

// Broadcast pushes a frame to every client. Drops frames for slow clients.
func (h *clientHub) Broadcast(frame []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.clients {
		select {
		case ch <- frame:
		default:
		}
	}
}
