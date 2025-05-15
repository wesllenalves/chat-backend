package ws

import (
	"log"
	"net/http"

	"chat-backend/internal/auth"
	"chat-backend/internal/presence"

	"github.com/gorilla/websocket"
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

	client := presence.NewClient(userID, conn)
	presence.Register(client)
	go client.Listen()
}
