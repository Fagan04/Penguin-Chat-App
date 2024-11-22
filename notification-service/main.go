package main

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/handlers"
	"github.com/Fagan04/Penguin-Chat-App/notification-service/repository"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "../database/notifications.db")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			message TEXT NOT NULL,
			is_new BOOLEAN NOT NULL,
			timestamp TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Pass the shared db instance to the repository
	repo := repository.NewRepository(db)
	r := mux.NewRouter()
	r.HandleFunc("/notifications/{user_id}", handlers.GetNewMessageHandler(repo)).Methods("GET")

	log.Println("Notification Service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
