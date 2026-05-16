package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

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

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(firstName, lastName, email string) (User, error) {
	newUser :=
		User{
			ID:        uuid.New(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			CreatedAt: time.Now(),
		}
	return s.repo.Create(newUser)
}

func (s *UserService) FindByID(id uuid.UUID) (User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Update(id uuid.UUID, user User) (User, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, err
	}

	existing.FirstName = user.FirstName
	existing.LastName = user.LastName
	existing.Email = user.Email

	return s.repo.Update(id, existing)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
