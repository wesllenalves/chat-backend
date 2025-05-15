package ws

import (
	"log"
	"net/http"

	"chat-backend/internal/auth"
	"chat-backend/internal/presence"
	"chat-backend/internal/chat"

	"github.com/gorilla/websocket"
	"database/sql"
	"fmt"

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
	db, err := GetDatabaseConnection() // Assuming GetDatabaseConnection() returns *sql.DB and error
	if err != nil {
		log.Println("Failed to get database connection:", err)
		return
	}
	store := chat.NewStore(db) // Pass the *sql.DB instance to NewStore
	client := presence.NewClient(userID, conn, store)
	presence.Register(client)
	go client.Listen()
}

func GetDatabaseConnection() (*sql.DB, error) {
	connStr := "user=username dbname=chatdb sslmode=disable" // Replace with your actual connection string
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
