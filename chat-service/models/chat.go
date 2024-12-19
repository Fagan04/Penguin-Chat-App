package models

import (
	"github.com/Fagan04/Penguin-Chat-App/user-service/models"
	"net/http"
	"time"
)

type Chat struct {
	ChatID   int    `json:"id"`
	ChatName string `json:"chat_name"`
	OwnerID  int    `json:"owner_id"`
}

type ChatMember struct {
	ChatMemberID int       `json:"chat_member_id"`
	ChatID       int       `json:"chat_id"`
	UserID       int       `json:"user_id"`
	JoinedAt     time.Time `json:"joined_at"`
}

type ChatMemberWithUsername struct {
	ChatMemberID int
	ChatID       int
	UserID       int
	JoinedAt     string
	Username     string
}

type ChatMessage struct {
	MessageID   int    `json:"message_id"`
	ChatID      int    `json:"chat_id"`
	UserID      int    `json:"user_id"`
	MessageText string `json:"message_text"`
	SentAt      string `json:"sent_at"`
}

type ChatMessageWithUsername struct {
	MessageID   int    `json:"message_id"`
	ChatID      int    `json:"chat_id"`
	UserID      int    `json:"user_id"`
	MessageText string `json:"message_text"`
	SentAt      string `json:"sent_at"`
	Username    string `json:"username"`
}

type RegisterChatPayload struct {
	ChatName string `json:"chat_name"`
}

type ChatStore interface {
	GetChatByName(chatName string) (*Chat, error)
	CreateChat(Chat) error
	GetUserChats(userID int) ([]Chat, error)
	SendMessage(message ChatMessage) error
	AddUserToChat(chatID, userID int) error
	GetChatByID(chatID int) (*Chat, error)
	GetChatMembers(chatID int) ([]ChatMember, error)
	GetMessagesByChats(userID int) (map[int][]ChatMessage, error)
	GetUserIDByUsername(username string) (int, error)
	ExtractUserIDFromToken(r *http.Request) (int, error)
	GetAllUsers() ([]models.User, error)
}