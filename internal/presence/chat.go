package presence

import (
	"encoding/json"
	"log"
	"strconv"

	"chat-backend/internal/chat"
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

func StartGroupChatSubscriber(store *chat.Store) {
	go func() {
		log.Println("📡 Listening on Redis Pub/Sub: group:*")
		pubsub := redisdb.GetClient().PSubscribe(redisdb.Ctx, "group:*")
		ch := pubsub.Channel()
		for msg := range ch {
			groupIDstr := msg.Channel[len("group:"):]
			var payload ChatPayload
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Println("Erro ao parsear payload de grupo:", err)
				continue
			}
			groupID, err := strconv.Atoi(groupIDstr)
			if err != nil {
				log.Println("ID de grupo inválido:", groupIDstr)
				continue
			}
			memberIDs, err := store.GetGroupMembers(groupID)
			if err != nil {
				log.Printf("Erro ao buscar membros do grupo %d: %v", groupID, err)
				continue
			}
			for _, userID := range memberIDs {
				if value, ok := Clients.Load(userID); ok {
					client := value.(*Client)
					client.Conn.WriteJSON(map[string]interface{}{
						"type":    "group",
						"groupId": groupID,
						"from":    payload.From,
						"message": payload.Message,
					})
				}
			}
		}
	}()
}

func GetGroupMembersFromCacheOrDB(store *chat.Store, groupID string) []string {
	id, err := strconv.Atoi(groupID)
	if err != nil {
		return []string{}
	}
	members, err := store.GetGroupMembers(id)
	if err != nil {
		log.Printf("Erro ao buscar membros do grupo %s: %v", groupID, err)
		return []string{}
	}
	return members
}
