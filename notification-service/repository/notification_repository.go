package repository

import (
	"database/sql"
	"github.com/fagan04/penguin-chat-app/notification-service/model"
	_ "github.com/mattn/go-sqlite3"
)

func FetchNewMessages(userID string) ([]model.Notification, error) {
	db, err := sql.Open("sqlite3", "./notifications.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM notifications WHERE user_id=? AND is_new = 1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notif model.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Message, &notif.IsNew, &notif.Timestamp); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	_, err = db.Exec("UPDATE notifications SET is_new = 0 WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
