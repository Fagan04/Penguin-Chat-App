package main

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/user-service/auth"
	"github.com/Fagan04/Penguin-Chat-App/user-service/handlers"
	"github.com/Fagan04/Penguin-Chat-App/user-service/repository"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "../database/user.db")

	if err != nil {
		log.Printf("Error: %v", err)
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
		log.Printf("failed to create table: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("failed to close db: %v", err)
			return
		}
	}(db)

	userRepo := &repository.UserRepository{DB: db}
	userHandler := &handlers.UserHandler{Repo: userRepo}

	r := mux.NewRouter()

	r.HandleFunc("/login", userHandler.LoginUser).Methods("POST")
	r.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with your frontend URL for stricter control
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		r.ServeHTTP(w, req)

	})

	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(auth.JWTMiddleware)

	log.Println("User services is running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))

}
