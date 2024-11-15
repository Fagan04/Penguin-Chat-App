package handlers

import (
	"encoding/json"
	"github.com/fagan04/penguin-chat-app/notification-service/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func GetNewMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	messages, err := repository.FetchNewMessages(userID)
	if err != nil {
		http.Error(w, "Failed to fetch new messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
