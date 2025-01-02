package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Error encoding JSON: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, ApiResponse{Message: message, Status: "error"})
}
