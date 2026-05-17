package repository

import (
	"api/internal/domain"
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

func (r *InMemoryUserRepo) Create(user domain.User) (domain.User, error) {
	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepo) FindByID(id uuid.UUID) (domain.User, error) {
	user, ok := r.users[id]
	if !ok {

		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}
func (r *InMemoryUserRepo) Update(id uuid.UUID, user domain.User) (domain.User, error) {
	_, ok := r.users[id]
	if !ok {

		return domain.User{}, ErrUserNotFound
	}
	r.users[id] = user
	return user, nil
}
func (r *InMemoryUserRepo) Delete(id uuid.UUID) error {
	_, ok := r.users[id]
	if !ok {

		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil

}

type InMemoryUserRepo struct {
	users map[uuid.UUID]domain.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[uuid.UUID]domain.User),
	}
}
