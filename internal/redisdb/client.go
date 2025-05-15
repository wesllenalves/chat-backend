package redisdb

import (
	"chat-backend/internal/chat"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type ChatContext struct {
	UserID string
	Store  *chat.Store
}

// Defina ChatPayload localmente para evitar o ciclo de importação
type ChatPayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

var (
	Ctx = context.Background()
	rdb *redis.Client
)

func Init() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("❌ Redis connection error: %v", err)
	}
	log.Println("✅ Connected to Redis")
}

func GetClient() *redis.Client {
	return rdb
}

func PublishMessage(c *ChatContext, to string, content string, okTo bool, okContent bool) {
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
		err := GetClient().Publish(Ctx, channel, payloadBytes).Err()
		if err != nil {
			log.Printf("❌ Erro ao publicar mensagem no Redis: %v", err)
		} else {
			log.Printf("🔔 Mensagem publicada no canal Redis: %s", channel)
		}
	}
}

func PublishGroupMessage(groupID int, from, content string) error {
	payload := map[string]interface{}{
		"from":    from,
		"message": content,
	}
	payloadBytes, _ := json.Marshal(payload)
	channel := fmt.Sprintf("group:%d", groupID)
	return GetClient().Publish(Ctx, channel, payloadBytes).Err()
}
