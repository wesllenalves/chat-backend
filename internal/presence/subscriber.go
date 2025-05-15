package presence

import (
	"log"
	"strings"

	"chat-backend/internal/redisdb"
)

func StartPresenceSubscriber() {
	pubsub := redisdb.GetClient().Subscribe(redisdb.Ctx, "presence-updates")

	go func() {
		log.Println("📡 Listening on Redis Pub/Sub: presence-updates")

		for msg := range pubsub.Channel() {
			parts := strings.Split(msg.Payload, " ")
			if len(parts) == 2 {
				log.Printf("[🔔] PRESENCE EVENT: %s", msg.Payload)
				broadcastAllClients()
			}
		}
	}()
}
