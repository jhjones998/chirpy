package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type chirpError struct {
	Error string `json:"error"`
}

type standardResponse struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(chirpError{Error: msg})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(payload)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
