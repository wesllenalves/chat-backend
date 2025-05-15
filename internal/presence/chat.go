package presence

import (
	"encoding/json"
	"log"

	"chat-backend/internal/redisdb"
)

type ChatPayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func StartChatSubscriber() {
	go func() {
		log.Println("📡 Listening on Redis Pub/Sub: chat:*")

		pubsub := redisdb.GetClient().PSubscribe(redisdb.Ctx, "chat:*")
		ch := pubsub.Channel()

		for msg := range ch {
			var payload ChatPayload
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Println("❌ Erro ao parsear chat payload:", err)
				continue
			}

			userID := msg.Channel[len("chat:"):] // extrai o ID do destinatário

			if client, ok := clients[userID]; ok {
				client.Conn.WriteJSON(map[string]interface{}{
					"type":    "chat",
					"from":    payload.From,
					"message": payload.Message,
				})
			}
		}
	}()
}
