package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func UserPost(w http.ResponseWriter, r *http.Request) {
	var bodyParsed CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&bodyParsed)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Succesfully post user :%v\n", bodyParsed)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bodyParsed); err != nil {
		log.Printf("encode failed: %v", err)
	}
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Succesfully get user with id:%v\n", vars)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vars); err != nil {
		log.Printf("encode failed: %v", err)
	}
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Succesfully delete user with id:%v\n", vars)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vars); err != nil {
		log.Printf("encode failed: %v", err)
	}
}

func UserPatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("We handle Pacho !!!!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
