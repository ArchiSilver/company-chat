package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"company-chat/internal/metrics"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	ChatID string
	// ограничение скорости (rate limiting)
	WindowStart time.Time
	MsgCount    int
}

type Hub struct {
	mu sync.RWMutex
	// chatID -> клиенты
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*Client]struct{})}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[c.ChatID]; !ok {
		h.rooms[c.ChatID] = make(map[*Client]struct{})
	}
	h.rooms[c.ChatID][c] = struct{}{}
	// обновляем метрики
	metrics.ActiveWSConnections.Inc()
	if metrics.WSByChatEnabled() {
		metrics.ActiveWSConnectionsByChat.WithLabelValues(c.ChatID).Inc()
	}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.rooms[c.ChatID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.rooms, c.ChatID)
		}
	}
	// обновляем метрики
	metrics.ActiveWSConnections.Dec()
	if metrics.WSByChatEnabled() {
		metrics.ActiveWSConnectionsByChat.WithLabelValues(c.ChatID).Dec()
	}
}

func (h *Hub) Broadcast(chatID string, msg []byte, except *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.rooms[chatID]; ok {
		for cl := range clients {
			if cl == except {
				continue
			}
			select {
			case cl.Send <- msg:
			default:
				// сброс (если канал переполнен)
			}
		}
	}
}
