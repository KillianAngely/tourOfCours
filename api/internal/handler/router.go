package handler

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	userRouter := r.PathPrefix("/user").Subrouter()
	// "/user/"
	userRouter.HandleFunc("", UserPost).Methods("POST")
	userRouter.HandleFunc("/{id}", UserGet).Methods("GET")
	userRouter.HandleFunc("/{id}", UserPatch).Methods("PATCH")
	userRouter.HandleFunc("/{id}", UserDelete).Methods("DELETE")

	r.HandleFunc("/", HomeHello).Methods("GET")

	return r
}
