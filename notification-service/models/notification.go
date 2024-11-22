package models

type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	IsNew     bool   `json:"is_new"`
	Timestamp string `json:"timestamp"`
}
