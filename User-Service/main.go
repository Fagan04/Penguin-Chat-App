package main

import (
	"database/sql"
	"github.com/fagan04/penguin-chat-app/user-service/auth"
	"github.com/fagan04/penguin-chat-app/user-service/handlers"
	"github.com/fagan04/penguin-chat-app/user-service/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite", "./database/user.db")
	if err != nil {
		log.Fatalf("failed to connect to SQLite database: %v", err)
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
