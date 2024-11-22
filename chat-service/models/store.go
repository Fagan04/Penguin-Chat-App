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

func (c *Store) GetChatByID(chatID int) (*Chat, error) {
	// Log that we're trying to get a chat by ID
	fmt.Printf("Attempting to get chat by ID: %d\n", chatID)

	rows, err := c.db.Query("SELECT chat_id, chat_name FROM chats WHERE chat_id = ?", chatID)
	if err != nil {
		// Log query error
		return nil, fmt.Errorf("failed to query chat by ID: %w", err)
	}
	defer rows.Close()

	// Log whether any rows were returned
	if !rows.Next() {
		fmt.Println("No rows found for chat_id:", chatID)
		return nil, fmt.Errorf("chat not found")
	}

	// Scan the result and log the values
	chat := new(Chat)
	err = rows.Scan(&chat.ChatID, &chat.ChatName)
	if err != nil {
		// Log scanning error
		return nil, fmt.Errorf("failed to scan chat row: %w", err)
	}

	// Log the values to verify
	fmt.Printf("Found chat: ID=%d, Name=%s\n", chat.ChatID, chat.ChatName)

	if chat.ChatName == "" {
		// Log if the chat_name is empty
		return nil, fmt.Errorf("chat name is empty")
	}

	return chat, nil
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

func (c *Store) AddUserToChat(userID, chatID int) error {
	joinedAt := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("Adding user %d to chat %d at %s\n", userID, chatID, joinedAt)
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

func (s *Store) GetChatMembers(chatID int) ([]ChatMember, error) {
	rows, err := s.db.Query("SELECT chat_member_id, chat_id, user_id, joined_at FROM chat_members WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chat members: %w", err)
	}
	defer rows.Close()

	var members []ChatMember
	for rows.Next() {
		var member ChatMember
		if err := rows.Scan(&member.ChatMemberID, &member.ChatID, &member.UserID, &member.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan chat member: %w", err)
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return members, nil
}
