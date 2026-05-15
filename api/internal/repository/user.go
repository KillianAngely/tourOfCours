package repository

import (
	"api/internal/service"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user service.User) (service.User, error)
	FindByID(id uuid.UUID) (service.User, error)
	Update(id uuid.UUID, user service.User) (service.User, error)
	Delete(id uuid.UUID) error
}
