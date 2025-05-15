package presence

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"chat-backend/internal/chat"
	"chat-backend/internal/redisdb"

	"github.com/gorilla/websocket"

	"sync"
)

var Clients sync.Map

type ChatPayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Store  *chat.Store
}

// Use the shared Clients variable from shared.go

func NewClient(userID string, conn *websocket.Conn, store *chat.Store) *Client {
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Store:  store,
	}
	Clients.Store(userID, client) // Store the client in the sync.Map
	return client
}

func (c *Client) Listen() {
	defer func() {
		c.Conn.Close()
		Clients.Delete(c.UserID)
	}()

	for {
		var msg map[string]interface{}
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("[📩] Message from %s: %v", c.UserID, msg)

		// --- Aqui entra o trecho para mensagem de grupo ---
		if groupIDf, ok := msg["group_id"].(float64); ok {
			content, okContent := msg["content"].(string)
			if okContent {
				groupID := int(groupIDf)
				if err := redisdb.PublishGroupMessage(groupID, c.UserID, content); err != nil {
					log.Printf("Erro ao publicar mensagem de grupo no Redis: %v", err)
				} else {
					log.Printf("🔔 Mensagem de grupo publicada no canal Redis: group:%d", groupID)
				}
			}
			continue
		}
		// --- Fim do trecho para mensagem de grupo ---

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

			// Publicar a mensagem no canal Redis
			payload := ChatPayload{
				From:    c.UserID,
				Message: content,
			}
			payloadBytes, _ := json.Marshal(payload)
			channel := "chat:" + to
			redisdb.GetClient().Publish(redisdb.Ctx, channel, payloadBytes)
			log.Printf("🔔 Mensagem publicada no canal Redis: %s", channel)

			// Notificar o cliente conectado
			Clients.Range(func(key, value interface{}) bool {
				log.Printf("📋 Cliente conectado: %v", key)
				return true
			})
			if value, ok := Clients.Load(to); ok {
				client := value.(*Client)
				log.Printf("🔔 Cliente encontrado: %s", to)
				client.Conn.WriteJSON(map[string]interface{}{
					"type":    "chat",
					"from":    payload.From,
					"message": payload.Message,
				})
			} else {
				log.Printf("❌ Cliente %s não está conectado", to)
			}
		}
	}
}
