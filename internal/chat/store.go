package chat

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID        int64
	From      string
	To        string
	Content   string
	Timestamp time.Time
}

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) SaveMessage(ctx context.Context, msg Message) error {
	query := `INSERT INTO messages (sender, receiver, content, timestamp)
	          VALUES ($1, $2, $3, $4)`
	_, err := s.DB.ExecContext(ctx, query, msg.From, msg.To, msg.Content, msg.Timestamp)
	return err
}
