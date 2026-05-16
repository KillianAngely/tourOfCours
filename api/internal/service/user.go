package service

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRepository interface {
	Create(user User) (User, error)
	FindByID(id uuid.UUID) (User, error)
	Update(id uuid.UUID, user User) (User, error)
	Delete(id uuid.UUID) error
}
