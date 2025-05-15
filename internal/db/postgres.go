package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "postgres"
	}
	dsn := fmt.Sprintf("postgres://chat:chatpass@%s:5432/chatdb?sslmode=disable", host)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ DB Connect error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("❌ DB Ping error:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db
}
