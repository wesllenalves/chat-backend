package presence

import (
	"log"
	"sync"

	"chat-backend/internal/redisdb"
)

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

func Register(c *Client) {
	mu.Lock()
	defer mu.Unlock()

	clients[c.UserID] = c
	log.Printf("✅ Cliente %s registrado no sistema de presença", c.UserID)

	redisdb.GetClient().SAdd(redisdb.Ctx, "online_users", c.UserID)
	redisdb.GetClient().Publish(redisdb.Ctx, "presence-updates", c.UserID+" online")
}

func Unregister(c *Client) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, c.UserID)
	log.Printf("❌ Cliente %s removido do sistema de presença", c.UserID)

	c.Conn.Close()
	redisdb.GetClient().SRem(redisdb.Ctx, "online_users", c.UserID)
	redisdb.GetClient().Publish(redisdb.Ctx, "presence-updates", c.UserID+" offline")
}

func broadcastAllClients() {
	online, err := redisdb.GetClient().SMembers(redisdb.Ctx, "online_users").Result()
	if err != nil {
		log.Println("Erro ao recuperar presença do Redis:", err)
		return
	}

	for _, client := range clients {
		client.Conn.WriteJSON(map[string]interface{}{
			"type":   "presence",
			"online": online,
		})
	}
}
