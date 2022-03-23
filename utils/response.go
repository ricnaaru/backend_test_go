package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SuccessResponse struct {
	Timestamp string                 `json:"timestamp"`
	Success   bool                   `json:"success"`
	Response  map[string]interface{} `json:"response"`
}

type FailedResponse struct {
	Timestamp string                 `json:"timestamp"`
	Success   bool                   `json:"success"`
	Error     map[string]interface{} `json:"error"`
}

func ReturnFailure(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	s := FailedResponse{Success: false, Timestamp: time.Now().Format(time.RFC3339), Error: map[string]interface{}{
		"message": message,
	}}

	jsonResp, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}

func ReturnUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	s := FailedResponse{Success: false, Timestamp: time.Now().Format(time.RFC3339), Error: map[string]interface{}{
		"message": message,
	}}

	jsonResp, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}

func ReturnSuccess(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	s := SuccessResponse{Success: false, Timestamp: time.Now().Format(time.RFC3339), Response: response}

	jsonResp, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}
