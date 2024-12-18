package repository

import (
	"database/sql"
	"log"
)

type ChatRepository struct {
	DB *sql.DB
}

func (repo *ChatRepository) CreateMessage(userID, message, timestamp string) error {
	_, err := repo.DB.Exec(`
		INSERT INTO messages (user_id, message, timestamp) 
		VALUES (?, ?, ?)
	`, userID, message, timestamp)
	if err != nil {
		log.Printf("Error inserting message: %v", err)
		return err
	}
	return nil
}

func (repo *ChatRepository) GetMessages(userID string) ([]Message, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_id, message, timestamp 
		FROM messages WHERE user_id = ?`, userID)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Message, &msg.Timestamp); err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return messages, nil
}

type Message struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
