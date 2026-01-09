package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"company-chat/internal/domain"
	"company-chat/internal/logging"
	"company-chat/internal/metrics"
	"company-chat/internal/repository"
)

type incomingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// HandleWS выполняет апгрейд соединения, аутентификацию и запускает циклы чтения/записи
func HandleWS(h *Hub, dbRepo *repository.DB, msgRepo *repository.MessageRepository, jwtSecret string, conn *websocket.Conn, userID, chatID string) {
	// настройки heartbeat и ограничения скорости
	const (
		pongWait   = 60 * time.Second
		pingPeriod = (pongWait * 9) / 10
		writeWait  = 10 * time.Second
		maxMsgSize = 1024
		msgWindow  = 10 * time.Second
		maxMsgs    = 20
	)

	client := &Client{Conn: conn, Send: make(chan []byte, 256), UserID: userID, ChatID: chatID}
	// поля лимитера сообщений для каждого клиента
	client.WindowStart = time.Now()
	client.MsgCount = 0

	h.Register(client)

	client.Conn.SetReadLimit(maxMsgSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// writer: отправка сообщений и периодические ping
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			client.Conn.Close()
		}()
		for {
			select {
			case b, ok := <-client.Send:
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// канал закрыт
					client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				if err := client.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
					logging.L.Warnf("ws write error: %v", err)
					return
				}
			case <-ticker.C:
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					logging.L.Warnf("ws ping error: %v", err)
					return
				}
			}
		}
	}()

	// reader: обработка входящих сообщений
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			logging.L.Infof("ws read error: %v", err)
			break
		}
		var im incomingMessage
		if err := json.Unmarshal(data, &im); err != nil {
			logging.L.Warnf("invalid ws payload: %v", err)
			continue
		}
		if im.Type == "message" {
			// простое ограничение скорости на клиента
			now := time.Now()
			if now.Sub(client.WindowStart) > msgWindow {
				client.WindowStart = now
				client.MsgCount = 0
			}
			client.MsgCount++
			if client.MsgCount > maxMsgs {
				logging.L.Warnf("rate limit exceeded for user %s", userID)
				continue
			}

			// сохранить сообщение в БД
			m := &domain.Message{
				ChatID:   chatID,
				SenderID: userID,
				Content:  im.Content,
			}
			if err := msgRepo.Create(context.Background(), m); err != nil {
				logging.L.Warnf("could not save message: %v", err)
				// инкрементировать метрику ошибок
				metrics.ErrorsTotal.WithLabelValues("ws_save").Inc()
				continue
			}
			// метрики: счётчик и размер
			metrics.MessagesSent.WithLabelValues("ws").Inc()
			metrics.MessageSize.Observe(float64(len(m.Content)))
			// рассылка сохранённого сообщения в формате JSON (включая id и created_at)
			b, err := json.Marshal(m)
			if err != nil {
				logging.L.Warnf("could not marshal message: %v", err)
				continue
			}
			h.Broadcast(chatID, b, client)
		}
	}

	h.Unregister(client)
}
