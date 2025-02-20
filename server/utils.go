package server

import (
	"encoding/json"
	"net/http"
)

// Success returns a 2xx status code on handlers
func Success(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}

// Error returns a 4xx status code on handlers for bad requests
func Error(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}
