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

func (s *Store) SaveGroupMessage(ctx context.Context, groupID int, sender, content string) error {
	_, err := s.DB.ExecContext(ctx,
		"INSERT INTO group_messages (group_id, sender, content) VALUES ($1, $2, $3)",
		groupID, sender, content)
	return err
}

func (s *Store) GetGroupMembers(groupID int) ([]string, error) {
	rows, err := s.DB.Query("SELECT user_id FROM group_members WHERE group_id = $1", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, nil
}
