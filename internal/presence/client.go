package presence

import (
	"context"
	"log"
	"time"

	"chat-backend/internal/chat"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Store  *chat.Store
}

func NewClient(userID string, conn *websocket.Conn, store *chat.Store) *Client {
	return &Client{UserID: userID, Conn: conn, Store: store}
}

func (c *Client) Listen() {
	defer func() {
		c.Conn.Close()
		Unregister(c)
	}()

	for {
		var msg map[string]interface{}
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("[📩] Message from %s: %v", c.UserID, msg)

		// Processar e salvar a mensagem
		to, okTo := msg["to"].(string)
		content, okContent := msg["content"].(string)
		if okTo && okContent {
			message := chat.Message{
				From:      c.UserID,
				To:        to,
				Content:   content,
				Timestamp: time.Now(),
			}
			if err := c.Store.SaveMessage(context.Background(), message); err != nil {
				log.Println("Erro ao salvar mensagem:", err)
			}
		}
	}
}
