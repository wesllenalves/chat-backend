package main

import (
	"log"
	"net/http"
	"os"

	"chat-backend/internal/auth"
	"chat-backend/internal/chat"
	"chat-backend/internal/db"
	"chat-backend/internal/presence"
	"chat-backend/internal/redisdb"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Arquivo .env não encontrado, usando variáveis de ambiente padrão")
	}

	// 1. Inicializar Redis
	log.Println("Inicializando Redis...")
	redisdb.Init()

	// 2. Conectar ao PostgreSQL e criar Store
	log.Println("Conectando ao PostgreSQL...")
	sqlDB := db.Connect()
	store := chat.NewStore(sqlDB)
	log.Println("Conectado ao PostgreSQL e Store criada.")

	// 3. Iniciar Pub/Sub para presença
	log.Println("Iniciando o assinante de presença...")
	presence.StartPresenceSubscriber()
	log.Println("Assinante de presença iniciado.")

	// 4. Endpoint para gerar JWT (GET /token?user=joao)
	http.HandleFunc("/token", auth.GenerateTokenHandler)

	// 5. Endpoint WebSocket para troca de mensagens
	http.HandleFunc("/ws", auth.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth.UserIDKey).(string)
		if !ok {
			log.Println("Erro ao obter userID do contexto")
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		log.Printf("Usuário %s tentando se conectar via WebSocket", userID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erro ao fazer upgrade para WebSocket:", err)
			return
		}

		log.Printf("Conexão WebSocket estabelecida para usuário %s", userID)

		client := presence.NewClient(userID, conn, store)
		log.Printf("Cliente %s criado", userID)

		presence.Register(client)
		log.Printf("Cliente %s registrado no sistema de presença", userID)

		go func() {
			defer func() {
				log.Printf("Cliente %s desconectado", userID)
				presence.Unregister(client)
			}()
			client.Listen()
		}()
	}))

	// 6. Iniciar o servidor
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
