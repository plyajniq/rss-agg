package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

// write 5XX error to response
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	RespondWithJSON(w, code, ErrResponse{Error: msg})
}

// write json to response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
