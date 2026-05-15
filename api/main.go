package main

import (
	"encoding/json"
	"fmt"
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
	http.ListenAndServe(":8000", r)
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
		fmt.Fprintf(w, "Bad Request")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succesfully create user with body :%v\n", bodyParsed)
}

func UserHandlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succesfully get user :%v\n", vars["id"])
}

func UserHandlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succesfully delete user with id :%v\n", vars["id"])
}

func UserHandlerPatch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "We handle Pacho !!!!")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
