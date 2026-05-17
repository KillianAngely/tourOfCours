package handler

import (
	"github.com/gorilla/mux"
)

func NewRouter(h *UserHandler) *mux.Router {
	r := mux.NewRouter()
	userRouter := r.PathPrefix("/user").Subrouter()
	// "/user/"
	userRouter.HandleFunc("", h.Post).Methods("POST")
	userRouter.HandleFunc("/{id}", h.Get).Methods("GET")
	userRouter.HandleFunc("/{id}", h.Patch).Methods("PATCH")
	userRouter.HandleFunc("/{id}", h.Delete).Methods("DELETE")

	r.HandleFunc("/", HealthHandler).Methods("GET")

	return r
}
