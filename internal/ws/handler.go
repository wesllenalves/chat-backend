package ws

import (
	"log"
	"net/http"

	"chat-backend/internal/auth"
	"chat-backend/internal/chat"
	"chat-backend/internal/presence"
	"chat-backend/internal/redisdb"

	"database/sql"
	"fmt"

	"github.com/gorilla/websocket"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	userID := r.Context().Value(auth.UserIDKey).(string)
	db, err := GetDatabaseConnection()
	if err != nil {
		log.Println("Failed to get database connection:", err)
		return
	}
	store := chat.NewStore(db)
	client := presence.NewClient(userID, conn, store)
	presence.Register(client)

	go func() {
		defer func() {
			conn.Close()
			presence.Unregister(client)
		}()
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("Read error:", err)
				break
			}
			// Mensagem de grupo
			if groupIDf, ok := msg["group_id"].(float64); ok {
				content, okContent := msg["content"].(string)
				if okContent {
					groupID := int(groupIDf)
					if err := redisdb.PublishGroupMessage(groupID, userID, content); err != nil {
						log.Printf("Erro ao publicar mensagem de grupo no Redis: %v", err)
					} else {
						log.Printf("🔔 Mensagem de grupo publicada no canal Redis: group:%d", groupID)
					}
				}
			}
			// Aqui você pode tratar outras mensagens (ex: chat 1:1)
		}
	}()
}

func GetDatabaseConnection() (*sql.DB, error) {
	connStr := "user=username dbname=chatdb sslmode=disable" // Replace with your actual connection string
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
