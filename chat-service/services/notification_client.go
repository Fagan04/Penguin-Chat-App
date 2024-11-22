package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type NotificationService struct {
	BaseURL string
}

func NewNotificationService(baseURL string) *NotificationService {
	return &NotificationService{BaseURL: baseURL}
}

func (ns *NotificationService) SendNotification(userID int, message string) error {
	payload := map[string]interface{}{
		"user_id":   userID,
		"message":   message,
		"is_new":    true,
		"timestamp": fmt.Sprintf("%v", time.Now().Format(time.RFC3339)),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}
	log.Printf("Sending notification to user %d: %s", userID, message)
	resp, err := http.Post(fmt.Sprintf("%s/addNotification", ns.BaseURL), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send notification request: %w", err)
	}
	log.Printf("Notification response status: %s", resp.Status)

	defer resp.Body.Close()

	log.Printf("Received response from notification service: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification services responded with status: %s", resp.Status)
	}

	return nil
}
