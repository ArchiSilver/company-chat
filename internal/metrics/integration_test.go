package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"company-chat/internal/metrics"
	"company-chat/internal/ws"
)

func TestHTTPAndWSMetricsIntegration(t *testing.T) {
	hub := ws.NewHub()

	mux := http.NewServeMux()
	// endpoint метрик
	mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// повторное использование promhttp.Handler потребовало бы импорта promhttp; держим просто
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// тестовый REST-эндпоинт, который записывает метрику сообщения
	mux.HandleFunc("/test-message", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		// симулируем создание сообщения
		metrics.MessagesSent.WithLabelValues("test").Inc()
		metrics.MessageSize.Observe(float64(len(body.Content)))
		w.WriteHeader(http.StatusCreated)
	})

	// тестовый WS-эндпоинт: регистрирует/удаляет клиента в хабе
	mux.HandleFunc("/test-ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade failed", http.StatusInternalServerError)
			return
		}
		client := &ws.Client{Conn: conn, Send: make(chan []byte, 1), UserID: "u1", ChatID: "chat1"}
		hub.Register(client)
		// эхо-луп: прочитать одно сообщение и закрыть
		_, _, err = conn.ReadMessage()
		if err == nil {
			conn.WriteMessage(websocket.TextMessage, []byte("ok"))
		}
		hub.Unregister(client)
		conn.Close()
	})

	server := httptest.NewServer(metrics.Middleware(mux))
	defer server.Close()

	// начальные значения метрик
	beforeMsgs := testutil.ToFloat64(metrics.MessagesSent.WithLabelValues("test"))
	beforeActive := testutil.ToFloat64(metrics.ActiveWSConnections)

	// POST a message
	client := &http.Client{Timeout: 2 * time.Second}
	reqBody := `{"content":"hello world"}`
	resp, err := client.Post(server.URL+"/test-message", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("post error: %v", err)
	}
	resp.Body.Close()

	// Дадим метрикам небольшое время для обновления
	time.Sleep(50 * time.Millisecond)

	afterMsgs := testutil.ToFloat64(metrics.MessagesSent.WithLabelValues("test"))
	if afterMsgs <= beforeMsgs {
		t.Fatalf("expected messages count to increase: before=%v after=%v", beforeMsgs, afterMsgs)
	}

	// WebSocket: connect
	wsURL := "ws" + server.URL[len("http"):]
	d := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	wsConn, _, err := d.Dial(wsURL+"/test-ws", nil)
	if err != nil {
		t.Fatalf("ws dial failed: %v", err)
	}
	// give hub time to register
	time.Sleep(50 * time.Millisecond)
	midActive := testutil.ToFloat64(metrics.ActiveWSConnections)
	if midActive <= beforeActive {
		t.Fatalf("expected active ws connections to increase: before=%v mid=%v", beforeActive, midActive)
	}
	// close ws
	wsConn.WriteMessage(websocket.TextMessage, []byte("hi"))
	wsConn.Close()
	time.Sleep(50 * time.Millisecond)
	afterActive := testutil.ToFloat64(metrics.ActiveWSConnections)
	if afterActive != beforeActive {
		t.Fatalf("expected active ws connections to return to before: before=%v after=%v", beforeActive, afterActive)
	}
}
