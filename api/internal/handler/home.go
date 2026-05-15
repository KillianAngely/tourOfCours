package handler

import (
	"log"
	"net/http"
)

func HomeHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello")); err != nil {
		log.Printf("hello handler failed:")
	}
}
