package models

import (
	"database/sql"
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/user-service/auth"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	db *sql.DB
}

func (s *Store) ExtractUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header not found")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid authorization header")
	}

	token := parts[1]

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.Atoi(claims.Id) // Convert back to int
	if err != nil {
		return 0, errors.New("invalid user_id in token")
	}
	return userID, nil
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}

}

func (c *Store) GetChatByName(chatName string) (*Chat, error) {
	rows, err := c.db.Query("SELECT chat_id, chat_name FROM chats WHERE chat_name = ?", chatName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("chat not found")
	}

	ch := new(Chat)
	err = rows.Scan(&ch.ChatID, &ch.ChatName)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *Store) GetChatByID(chatID int) (*Chat, error) {
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

	chat := new(Chat)
	err = rows.Scan(&chat.ChatID, &chat.ChatName)
	if err != nil {
		// Log scanning error
		return nil, fmt.Errorf("failed to scan chat row: %w", err)
	}

	fmt.Printf("Found chat: ID=%d, Name=%s\n", chat.ChatID, chat.ChatName)

	if chat.ChatName == "" {
		return nil, fmt.Errorf("chat name is empty")
	}

	return chat, nil
}

func (c *Store) CreateChat(chat Chat) error {
	query := "INSERT INTO chats (chat_name, owner_id) VALUES (?, ?)"

	result, err := c.db.Exec(query, chat.ChatName, chat.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	chat.ChatID = int(chatID)

	joinedAt := time.Now().Format("2006-01-02 15:04:05")
	memberQuery := "INSERT INTO chat_members (chat_id, user_id, joined_at) VALUES (?, ?, ?)"
	_, err = c.db.Exec(memberQuery, chat.ChatID, chat.OwnerID, joinedAt)
	if err != nil {
		return fmt.Errorf("failed to add chat owner as a member: %w", err)
	}

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

func (s *Store) GetUserChats(userID int) ([]Chat, error) {
	query := `
		SELECT c.chat_id, c.chat_name, c.owner_id
		FROM chats c
		JOIN chat_members cm ON c.chat_id = cm.chat_id
		WHERE cm.user_id = ?`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := []Chat{}
	for rows.Next() {
		chat := Chat{}
		if err := rows.Scan(&chat.ChatID, &chat.ChatName, &chat.OwnerID); err != nil {
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

func (s *Store) GetMessagesByChats(userID int) (map[int][]ChatMessage, error) {
	// SQL query to fetch messages for the user, grouped by chat_id
	query := `
		SELECT m.chat_id, m.message_id, m.user_id, m.message_text, m.sent_at
		FROM chat_messages m
		JOIN chat_members cm ON m.chat_id = cm.chat_id
		WHERE cm.user_id = ?  -- Only messages for this user
		ORDER BY m.chat_id, m.sent_at;
	`

	// Execute the query, passing the userID as a parameter
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	groupedMessages := make(map[int][]ChatMessage)

	for rows.Next() {
		var msg ChatMessage
		var chatID int

		err := rows.Scan(&chatID, &msg.MessageID, &msg.UserID, &msg.MessageText, &msg.SentAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		groupedMessages[chatID] = append(groupedMessages[chatID], msg)
	}

	// Return the map containing the messages grouped by chat ID
	return groupedMessages, nil
}

func (s *Store) GetUserIDByUsername(username string) (int, error) {
	var userID int
	query := `SELECT user_id FROM users WHERE username = ?`
	err := s.db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return userID, nil
}
