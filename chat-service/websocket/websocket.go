package websocket

import (
	"encoding/json"
	"github.com/Fagan04/Penguin-Chat-App/user-service/auth"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

type WebSocketServer struct {
	Clients     map[*websocket.Conn]bool
	ChatClients map[int]map[*websocket.Conn]bool
	Broadcast   chan []byte
	Register    chan *websocket.Conn
	Unregister  chan *websocket.Conn
	mu          sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	log.Printf("User %s authenticated", claims.Username)

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("WebSocket upgrade failed: %v", err)
		return
	}

	s.Register <- conn

	chatIDStr := r.URL.Query().Get("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		log.Fatalf("Error extracting ChatID from request: ", err)
		conn.Close()
		return
	}

	s.mu.Lock()
	if _, ok := s.ChatClients[chatID]; !ok {
		s.ChatClients[chatID] = make(map[*websocket.Conn]bool)
	}
	s.ChatClients[chatID][conn] = true
	s.mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			s.Unregister <- conn
			err := conn.Close()
			if err != nil {
				log.Printf("Close error: %v", err)
			}
			break
		}
		// Broadcast the received message
		s.Broadcast <- msg
	}
}

func (s *WebSocketServer) Start() {
	for {
		select {
		case conn := <-s.Register:
			s.mu.Lock()
			s.Clients[conn] = true
			s.mu.Unlock()
			log.Println("New WebSocket client connected")

		case conn := <-s.Unregister:
			s.mu.Lock()
			if _, ok := s.Clients[conn]; ok {
				delete(s.Clients, conn)
				err := conn.Close()
				if err != nil {
					log.Printf("WebSocket client unregister error: %v", err)
				}
				log.Println("WebSocket client disconnected")
			}
			s.mu.Unlock()

		case message := <-s.Broadcast:

			var msgData struct {
				ChatID int    `json:"chat_id"`
				Text   string `json:"text"`
			}

			if err := json.Unmarshal(message, &msgData); err != nil {
				log.Println("Error unmarshaling message: ", err)
				continue
			}

			s.mu.Lock()
			if chatClients, ok := s.ChatClients[msgData.ChatID]; ok {
				for client := range chatClients {
					err := client.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Printf("Broadcast error: %v", err)
						err := client.Close()
						if err != nil {
							log.Printf("Close error: %v", err)
						}
						delete(chatClients, client)
					}
				}
			}
			s.mu.Unlock()
		}
	}
}
