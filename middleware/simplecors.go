package middleware

import (
	"net/http"
)

type (
	SimpleCORS struct {
	}
)

func NewSimpleCORS() *SimpleCORS {
	return &SimpleCORS{}
}

func (m SimpleCORS) Handle(w http.ResponseWriter, r *http.Request) (err *error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// pre-flight
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	return
}
