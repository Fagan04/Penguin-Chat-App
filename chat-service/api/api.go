package api

import (
	"database/sql"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/handlers"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/models"
	"github.com/Fagan04/Penguin-Chat-App/chat-service/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr   string
	db     *sql.DB
	router *mux.Router
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api").Subrouter()

	chatStore := models.NewStore(s.db)
	chatHandler := handlers.NewChatHandler(chatStore)
	chatHandler.RegisterRoutes(subrouter)

	wsServer := websocket.NewWebSocketServer()
	go wsServer.Start()

	router.HandleFunc("/ws", wsServer.HandleConnections)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
