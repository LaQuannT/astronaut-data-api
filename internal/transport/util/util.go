package util

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
