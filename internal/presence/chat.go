package presence

import (
	// "chat-backend/internal/presence/client" // Removed as the package is not available
	"encoding/json"
	"log"

	"chat-backend/internal/redisdb"
)

func StartChatSubscriber() {
	go func() {
		log.Println("📡 Listening on Redis Pub/Sub: chat:*")

		pubsub := redisdb.GetClient().PSubscribe(redisdb.Ctx, "chat:*")
		ch := pubsub.Channel()

		for msg := range ch {
			log.Printf("🔔 Mensagem recebida no canal Redis: %s", msg.Channel)
			var payload ChatPayload
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				// Removed unused variable declaration
				continue
			}

			userID := msg.Channel[len("chat:"):] // extrai o ID do destinatário
			log.Printf("🔔 Mensagem destinada ao usuário: %s", userID)

			if client, ok := clients[userID]; ok {
				log.Printf("🔔 Enviando mensagem para o cliente conectado: %s", userID)
				client.Conn.WriteJSON(map[string]interface{}{
					"type":    "chat",
					"from":    payload.From,
					"message": payload.Message,
				})
			} else {
				log.Printf("❌ Cliente %s não está conectado", userID)
			}
		}
	}()
}
