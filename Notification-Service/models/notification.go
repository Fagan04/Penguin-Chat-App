package models

type Notification struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`   // ID of the user receiving the notification
	Message   string `json:"message"`   // Content of the notification
	IsNew     bool   `json:"is_new"`    // Flag indicating whether the message is new
	Timestamp string `json:"timestamp"` // Timestamp when the notification was created
}
