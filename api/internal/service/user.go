package service

import (
	"api/internal/domain"
	"api/internal/repository"

	"errors"
	"time"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
}

type UpdateUserInput struct {
	FirstName string
	LastName  string
	Email     string
}

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	FindByID(id uuid.UUID) (domain.User, error)
	Update(id uuid.UUID, user domain.User) (domain.User, error)
	Delete(id uuid.UUID) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(input CreateUserInput) (domain.User, error) {
	newUser :=
		domain.User{
			ID:        uuid.New(),
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			CreatedAt: time.Now(),
		}
	return s.repo.Create(newUser)
}

func (s *UserService) FindByID(id uuid.UUID) (domain.User, error) {
	user, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *UserService) Update(id uuid.UUID, input UpdateUserInput) (domain.User, error) {
	existing, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	existing.FirstName = input.FirstName
	existing.LastName = input.LastName
	existing.Email = input.Email

	return s.repo.Update(id, existing)
}

func (s *UserService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.ErrUserNotFound
	}
	return err
}
