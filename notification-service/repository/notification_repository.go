package repository

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/models"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (repo *NotificationRepository) FetchNewMessages(userID string) ([]models.Notification, error) {
	rows, err := repo.DB.Query("SELECT * FROM notifications WHERE user_id=? AND is_new = 1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Message, &notif.IsNew, &notif.Timestamp); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	_, err = repo.DB.Exec("UPDATE notifications SET is_new = 0 WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
