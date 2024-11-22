package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/models"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetNewMessageHandler(repo *repository.NotificationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]

		messages, err := repo.FetchNewMessages(userID)
		if err != nil {
			http.Error(w, "Failed to fetch new messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
			return
		}
	}
}

type NotificationService struct {
	repository *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repository: repo}
}

func (s *NotificationService) AddNotification(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Log the received notification
	log.Printf("Received notification: %v", notification)

	err := s.repository.SaveNotification(&notification)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save notification: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
