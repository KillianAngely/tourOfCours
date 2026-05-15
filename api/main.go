package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	userRouter := r.PathPrefix("/user").Subrouter()
	// "/user/"
	userRouter.HandleFunc("", UserHandlerPost).Methods("POST")
	userRouter.HandleFunc("/{id}", UserHandlerGet).Methods("GET")
	userRouter.HandleFunc("/{id}", UserHandlerPatch).Methods("PATCH")
	userRouter.HandleFunc("/{id}", UserHandlerDelete).Methods("DELETE")

	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

type UserBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func UserHandlerPost(w http.ResponseWriter, r *http.Request) {
	var bodyParsed UserBody
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

func UserHandlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Succesfully get user with id:%v\n", vars)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vars); err != nil {
		log.Printf("encode failed: %v", err)
	}
}

func UserHandlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Succesfully delete user with id:%v\n", vars)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vars); err != nil {
		log.Printf("encode failed: %v", err)
	}
}

func UserHandlerPatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("We handle Pacho !!!!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello")); err != nil {
		log.Printf("hello handler failed:")
	}
}
