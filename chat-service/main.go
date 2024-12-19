package main

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/handlers"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/models"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/websocket"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {

	db, err := sql.Open("sqlite3", "../database/chats.db")
	if err != nil {
		log.Fatal(err)
	}

	dbUser, err := sql.Open("sqlite3", "../database/user.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			chat_id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_name TEXT NOT NULL,
		    owner_id INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS chat_members (
		    chat_member_id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER,
			user_id INTEGER,
 			joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Use DATETIME or TIMESTAMP
    		FOREIGN KEY (chat_id) REFERENCES chats(chat_id),
    		FOREIGN KEY (user_id) REFERENCES users(user_id)	
		);
		CREATE TABLE IF NOT EXISTS chat_messages (
			message_id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER,
			user_id INTEGER,
			message_text TEXT,
			sent_at TEXT
		);

`)

	if err != nil {
		log.Printf("failed to create table: %v", err)
	}

	chatStore := models.NewStore(db, dbUser)
	chatHandler := handlers.NewChatHandler(chatStore)

	r := mux.NewRouter()
	chatHandler.RegisterRoutes(r)

	wsServer := websocket.NewWebSocketServer()
	go wsServer.Start()
	r.HandleFunc("/ws", wsServer.HandleConnections)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with your frontend URL for stricter control
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, chat_id")

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		r.ServeHTTP(w, req)

	})

	log.Println("Chat services is running on port 8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", handler))

}
