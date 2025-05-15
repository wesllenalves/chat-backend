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
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	database := os.Getenv("PGDATABASE")
	port := os.Getenv("PGPORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

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
