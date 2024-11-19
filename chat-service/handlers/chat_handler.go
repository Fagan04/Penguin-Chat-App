package handlers

import (
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/models"
	"github.com/Fagan04/Penguin-Chat-App/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	store models.ChatStore
}

func NewChatHandler(store models.ChatStore) *ChatHandler {
	return &ChatHandler{store: store}
}

func (c *ChatHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/createChat", c.HandlerChatCreation).Methods("POST")
	router.HandleFunc("/accessChat", c.HandlerChatAccess).Methods("GET")
	router.HandleFunc("/sendMessage", c.HandlerSendMessage).Methods("POST")
}

func (c *ChatHandler) HandlerChatCreation(w http.ResponseWriter, r *http.Request) {
	//get the json payload
	var payload models.RegisterChatPayload
	if err := utils.ParseJson(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// check if the chat exists
	_, err := c.store.GetChatByName(payload.ChatName)
	if err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("chat with name %s already exists", payload.ChatName))
		return
	}

	// if not
	err = c.store.CreateChat(models.Chat{
		ChatName: payload.ChatName,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	err = utils.WriteJson(w, http.StatusCreated, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

}

func (c *ChatHandler) HandlerChatAccess(w http.ResponseWriter, r *http.Request) {
	// Retrieve User ID from headers (set by the middleware)
	userIDStr := r.Header.Get("User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Fetch chats where the user is a member
	chats, err := c.store.GetUserChats(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with chat data
	err = utils.WriteJson(w, http.StatusOK, chats)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (c *ChatHandler) HandlerSendMessage(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

	//notify

	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("user not part of the chat"))
		return
	}

}