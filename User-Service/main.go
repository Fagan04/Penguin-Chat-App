package main

import (
	"database/sql"
	"github.com/fagan04/penguin-chat-app/user-service/auth"
	"github.com/fagan04/penguin-chat-app/user-service/handlers"
	"github.com/fagan04/penguin-chat-app/user-service/repository"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/user.db")
	if err != nil {
		log.Fatalf("failed to connect to SQLite database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL
    )
`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	defer db.Close()

	userRepo := &repository.UserRepository{DB: db}
	userHandler := &handlers.UserHandler{Repo: userRepo}

	r := mux.NewRouter()

	r.HandleFunc("/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")

	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(auth.JWTMiddleware)

	log.Println("User service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
