package repository

import (
	"database/sql"
	"fmt"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/models"
	"github.com/go-sql-driver/mysql"
	"log"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
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

func (repo *NotificationRepository) SaveNotification(notification *models.Notification) error {
	// Log the notification before saving
	log.Printf("Saving notification: %v", notification)

	_, err := repo.DB.Exec("INSERT INTO notifications (user_id, message, is_new, timestamp) VALUES (?, ?, ?, ?)",
		notification.UserID, notification.Message, notification.IsNew, notification.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	return nil
}
