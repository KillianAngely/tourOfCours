package handler

import (
	"io"
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, `{"alive": true}`)
	if err != nil {
		log.Printf("GET /health - error: %v", err)

	}
}
