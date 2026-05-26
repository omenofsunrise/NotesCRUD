package handlers

import (
	"encoding/json"
	"net/http"
)

func setErrorResponse(w http.ResponseWriter, status int, message string) {
	setResponse(w, status, map[string]string{"error": message})
}

func setResponse(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
