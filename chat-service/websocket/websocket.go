package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
		Clients:     make(map[*websocket.Conn]bool),
		ChatClients: make(map[int]map[*websocket.Conn]bool), // Add this
		Broadcast:   make(chan []byte),
		Register:    make(chan *websocket.Conn),
		Unregister:  make(chan *websocket.Conn),
	}

}
func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting to handle WebSocket connection from: %s", r.RemoteAddr)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	// Register the connection (e.g., add it to a default chat group)
	s.mu.Lock()
	defaultChatID := 0 // Use a default identifier for anonymous or general connections
	if _, ok := s.ChatClients[defaultChatID]; !ok {
		s.ChatClients[defaultChatID] = make(map[*websocket.Conn]bool)
	}
	s.ChatClients[defaultChatID][conn] = true
	s.mu.Unlock()

	s.Register <- conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			s.Unregister <- conn
			break
		}
		// Broadcast the received message to the default chat group
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
				Text string `json:"text"`
			}

			// Attempt to unmarshal as JSON
			if err := json.Unmarshal(message, &msgData); err != nil {
				log.Printf("Error unmarshaling message, treating as plain text: %v", err)
				msgData.Text = string(message) // Treat as plain text
			}

			log.Printf("Broadcasting message: %s", msgData.Text)

			s.mu.Lock()
			for client := range s.Clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(msgData.Text))
				if err != nil {
					log.Printf("Broadcast error: %v", err)
					client.Close()
					delete(s.Clients, client)
				}
			}
			s.mu.Unlock()

		}
	}
}
