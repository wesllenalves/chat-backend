package presence

import (
	"log"

	"chat-backend/internal/redisdb"
)

var clients = make(map[string]*Client)

func Register(c *Client) {
	clients[c.UserID] = c
	redisdb.GetClient().SAdd(redisdb.Ctx, "online_users", c.UserID)
	redisdb.GetClient().Publish(redisdb.Ctx, "presence-updates", c.UserID+" online")
}

func Unregister(c *Client) {
	delete(clients, c.UserID)
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
