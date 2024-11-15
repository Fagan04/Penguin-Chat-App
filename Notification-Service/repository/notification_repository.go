package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func FetchNewMessages(userID string) ([]map[string]interface{}, error) {
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

	var messages []map[string]interface{}
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, err
		}

		messages = append(messages, map[string]interface{}{
			"message": message,
		})
	}

	_, err = db.Exec("UPDATE notifications SET is_new = 0 WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
