package infrastructure

import (
	"chat_service/internal/domain"
	"database/sql"
	"fmt"
	"log"
)

type DBrepo struct {
	DB *sql.DB
}

func NewDBrepo(db *sql.DB) *DBrepo {
	return &DBrepo{DB: db}
}

func (db *DBrepo) SaveMessage(msg domain.Message) (domain.Message, error) {
	query := `INSERT INTO messages(user_id,username,content,role) 
	VALUES($1,$2,$3,$4)
	RETURNING id,created_at
	`
	err := db.DB.QueryRow(query, msg.UserID, msg.Username, msg.Content, msg.Role).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to save message: %w", err)
	}
	return msg, nil
}

func (db *DBrepo) MessageHistory(limit, offset int64) ([]domain.Message, error) {
	query := `
	SELECT id, user_id, username, content, created_at, role, is_anon
	FROM messages
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		return []domain.Message{}, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message

		err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Content, &msg.CreatedAt, &msg.Role, &msg.IsAnon)
		if err != nil {
			return []domain.Message{}, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (db *DBrepo) DeleteMessage(id int64) error {
	query := `DELETE FROM messages WHERE id = $1;`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("message not found:%v", err)
		return err
	}

	return nil
}
