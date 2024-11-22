package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	log.Println("Parsing JSON payload...")
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if err := WriteJson(w, status, map[string]string{"error": err.Error()}); err != nil {
		log.Printf("Failed to encode the response: %v", err)
	}
}
