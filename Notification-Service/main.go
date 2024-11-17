package main

import (
	"github.com/fagan04/penguin-chat-app/notification-service/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/notifications/{user_id}", handlers.GetNewMessageHandler).Methods("GET")

	log.Println("Notification Service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
