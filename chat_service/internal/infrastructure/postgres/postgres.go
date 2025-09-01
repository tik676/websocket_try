package postgres

import (
	"chat_service/internal/domain"
	"database/sql"
	"errors"
)

type Database struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *Database {
	return &Database{DB: db}
}

func (db *Database) SendMessage(message domain.Message) (domain.Message, error) {
	query := `INSERT INTO messages (user_id, user_name, role, is_anon, content)
          VALUES ($1, $2, $3, $4, $5)
          RETURNING id, created_at;`

	var msg domain.Message
	err := db.DB.QueryRow(
		query,
		message.UserID,
		message.UserName,
		message.Role,
		message.IsAnon,
		message.Content,
	).Scan(&msg.ID, &msg.CreatedAt)

	if err != nil {
		return domain.Message{}, errors.New("failed to send message")
	}

	msg.UserID = message.UserID
	msg.UserName = message.UserName
	msg.Role = message.Role
	msg.IsAnon = message.IsAnon
	msg.Content = message.Content

	msg.UserID = message.UserID
	msg.Content = message.Content

	return msg, nil
}

func (db *Database) DeleteMessage(id int64) error {
	query := `DELETE FROM messages WHERE id = $1;`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return errors.New("Failed message not exists")
	}
	return nil
}

func (db *Database) GetMessages(limit, offset int) ([]domain.Message, error) {
	query := `SELECT id, user_id, user_name, role, is_anon, content, created_at
          FROM messages
          ORDER BY created_at DESC
          LIMIT $1 OFFSET $2;`

	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		return []domain.Message{}, errors.New("Failed messages not exists")
	}

	var messages []domain.Message

	for rows.Next() {
		var msg domain.Message

		err := rows.Scan(&msg.ID, &msg.UserID, &msg.UserName, &msg.Role, &msg.IsAnon, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
