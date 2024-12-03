package handlers

import (
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/models"
	services "github.com/Fagan04/Penguin-Chat-App/chat-service/services"
	"github.com/Fagan04/Penguin-Chat-App/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	store               *models.Store
	notificationService *services.NotificationService
}

func NewChatHandler(store *models.Store, notificationService *services.NotificationService) *ChatHandler {
	return &ChatHandler{store: store, notificationService: notificationService}
}

func (c *ChatHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/createChat", c.HandlerChatCreation).Methods("POST")
	router.HandleFunc("/accessChat", c.HandlerChatAccess).Methods("GET")
	router.HandleFunc("/sendMessage", c.HandlerSendMessage).Methods("POST")
	router.HandleFunc("/addUserToChat", c.HandlerAddUserToChat).Methods("POST")
	router.HandleFunc("/getMessagesGroupedByChat", c.GetMessagesGroupedByChat).Methods("GET")
	router.HandleFunc("/getAllUsers", c.GetAllUsers).Methods("GET")
	router.HandleFunc("/getChatMembers", c.GetChatParticipants).Methods("GET")
}

func (c *ChatHandler) HandlerChatCreation(w http.ResponseWriter, r *http.Request) {
	userID, err := c.store.ExtractUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %w", err))
		return
	}

	var payload models.RegisterChatPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = c.store.GetChatByName(payload.ChatName)
	if err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("chat with name %s already exists", payload.ChatName))
		return
	}

	newChat := models.Chat{
		ChatName: payload.ChatName,
		OwnerID:  userID,
	}

	if err := c.store.CreateChat(newChat); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create chat: %w", err))
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "chat created successfully"})
}

func (c *ChatHandler) HandlerChatAccess(w http.ResponseWriter, r *http.Request) {
	userID, err := c.store.ExtractUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %w", err))
		return
	}

	chats, err := c.store.GetUserChats(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if len(chats) == 0 {
		utils.WriteJson(w, http.StatusOK, map[string]string{"message": "no chats found for the user"})
		return
	}

	err = utils.WriteJson(w, http.StatusOK, chats)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (c *ChatHandler) HandlerSendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := c.store.ExtractUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %w", err))
		return
	}

	var message models.ChatMessage
	if err := utils.ParseJson(r, &message); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = c.store.GetUserChats(message.ChatID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	message.UserID = userID
	message.SentAt = time.Now()

	err = c.store.SendMessage(message)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("About to send notifications")
	log.Printf("Message chat ID: %d", message.ChatID)
	chatMembers, err := c.store.GetChatMembers(message.ChatID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("Chat members: %+v", chatMembers)

	for _, member := range chatMembers {
		if member.UserID != userID {
			log.Println("Sending notification for chat message.")
			err := c.notificationService.SendNotification(member.UserID, fmt.Sprintf("New message in chat %d", message.ChatID))
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send notification: %w", err))
				return
			}
		}
	}

	// Respond with success
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "message sent successfully"})
}

func (c *ChatHandler) HandlerAddUserToChat(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ChatID   int    `json:"chat_id"`
		Username string `json:"username"` // Accept username instead of user ID
	}

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	chat, err := c.store.GetChatByID(payload.ChatID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	userID, err := c.store.ExtractUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %w", err))
		return
	}

	if chat.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you are not the owner of the chat"))
		return
	}

	targetUserID, err := c.store.GetUserIDByUsername(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found: %w", err))
		return
	}

	err = c.store.AddUserToChat(targetUserID, payload.ChatID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add user to chat: %w", err))
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "User added to chat successfully"})
}

func (c *ChatHandler) GetMessagesGroupedByChat(w http.ResponseWriter, r *http.Request) {
	userID, err := c.store.ExtractUserIDFromToken(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	groupedMessages, err := c.store.GetMessagesByChats(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, groupedMessages)
}

func (c *ChatHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := c.store.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, allUsers)
}

func (c *ChatHandler) GetChatParticipants(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.Header.Get("chat_id")
	if chatIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing chat ID"))
		return
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	participants, err := c.store.GetChatMembers(chatID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, participants)
}
