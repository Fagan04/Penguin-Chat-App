package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}

}

func (c *Store) GetChatByName(chatName string) (*Chat, error) {
	// Query the database for the chat by name
	rows, err := c.db.Query("SELECT chat_id, chat_name FROM chats WHERE chat_name = ?", chatName)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after we are done with them

	// Check if any rows are returned
	if !rows.Next() {
		return nil, fmt.Errorf("chat not found")
	}

	// Scan the result into a Chat struct
	ch := new(Chat)
	err = rows.Scan(&ch.ChatID, &ch.ChatName)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func scanRowIntoChat(row *sql.Rows) (*Chat, error) {
	chat := new(Chat)
	err := row.Scan(
		&chat.ChatID,
		&chat.ChatName,
	)
	if err != nil {
		return nil, err
	}

	return chat, err
}

func (c *Store) CreateChat(chat Chat) error {
	query := "INSERT INTO chats (chat_name) VALUES (?)"

	result, err := c.db.Exec(query, chat.ChatName)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	chat.ChatID = int(chatID)

	return nil
}

func (c *Store) AddUserToChat(chatID, userID int) error {
	joinedAt := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO chat_members (chat_id, user_id, joined_at) VALUES (?, ?, ?)"

	_, err := c.db.Exec(query, chatID, userID, joinedAt)
	if err != nil {
		return fmt.Errorf("failed to add user to chat: %w", err)
	}
	return nil
}

func (c *Store) GetUserChats(userID int) ([]Chat, error) {
	query := `
		SELECT c.chat_id, c.chat_name
		FROM chats c
		JOIN chat_members cm ON c.chat_id = cm.chat_id
		WHERE cm.user_id = ?`

	rows, err := c.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := []Chat{}
	for rows.Next() {
		chat := Chat{}
		if err := rows.Scan(&chat.ChatID, &chat.ChatName); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (c *Store) SendMessage(message ChatMessage) error {
	query := "INSERT INTO chat_messages (chat_id, user_id, message_text, sent_at) VALUES (?, ?, ?, ?)"

	_, err := c.db.Exec(query, message.ChatID, message.UserID, message.MessageText, message.SentAt)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
