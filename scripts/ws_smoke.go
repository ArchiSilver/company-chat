package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := flag.String("url", "ws://localhost:8080/ws", "ws url")
	token := flag.String("token", "", "access token")
	chat := flag.String("chat", "", "chat id")
	flag.Parse()
	if *token == "" || *chat == "" {
		log.Fatal("token and chat are required")
	}
	u := fmt.Sprintf("%s?token=%s&chat_id=%s", *url, *token, *chat)
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer c.Close()
	msg := `{"type":"message","content":"hello from go client"}`
	if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Fatalf("write: %v", err)
	}
	c.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, m, err := c.ReadMessage()
	if err != nil {
		log.Printf("read error (may be ok if no broadcast): %v", err)
	} else {
		log.Printf("recv: %s", string(m))
	}
}
