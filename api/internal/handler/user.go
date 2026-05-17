package handler

import (
	"api/internal/domain"
	"api/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users")
	var bodyParsed CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&bodyParsed)
	if err != nil {
		log.Printf("POST /users - bad request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := service.CreateUserInput{
		FirstName: bodyParsed.FirstName,
		LastName:  bodyParsed.LastName,
		Email:     bodyParsed.Email,
	}

	user, err := h.service.Create(input)
	if err != nil {
		log.Printf("POST /users - service error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("POST /users - created user %s", user.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("POST /users - encode failed: %v", err)
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /users/%s", mux.Vars(r)["id"])
	incomingId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("GET /users - invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.service.FindByID(incomingId)
	if errors.Is(err, domain.ErrUserNotFound) {
		log.Printf("GET /users/%s - not found", incomingId)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GET /users/%s - service error: %v", incomingId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("GET /users/%s - ok", incomingId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("GET /users/%s - encode failed: %v", incomingId, err)
	}
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE /users/%s", mux.Vars(r)["id"])
	incomingId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("DELETE /users - invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.Delete(incomingId)
	if errors.Is(err, domain.ErrUserNotFound) {
		log.Printf("DELETE /users/%s - not found", incomingId)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("DELETE /users/%s - service error: %v", incomingId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("DELETE /users/%s - deleted", incomingId)
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	log.Printf("PATCH /users/%s", mux.Vars(r)["id"])
	incomingId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("PATCH /users - invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var bodyParsed UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&bodyParsed)
	if err != nil {
		log.Printf("PATCH /users/%s - bad request: %v", incomingId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := service.UpdateUserInput{
		FirstName: bodyParsed.FirstName,
		LastName:  bodyParsed.LastName,
		Email:     bodyParsed.Email,
	}

	user, err := h.service.Update(incomingId, input)
	if errors.Is(err, domain.ErrUserNotFound) {
		log.Printf("PATCH /users/%s - not found", incomingId)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("PATCH /users/%s - service error: %v", incomingId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("PATCH /users/%s - updated", incomingId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("PATCH /users/%s - encode failed: %v", incomingId, err)
	}
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}
