package repository

import (
	"api/internal/service"

	"github.com/google/uuid"
)

func (r *InMemoryUserRepo) Create(user service.User) (service.User, error) {
	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepo) FindByID(id uuid.UUID) (service.User, error) {
	user, ok := r.users[id]
	if !ok {

		return service.User{}, service.ErrUserNotFound
	}
	return user, nil
}
func (r *InMemoryUserRepo) Update(id uuid.UUID, user service.User) (service.User, error) {
	_, ok := r.users[id]
	if !ok {

		return service.User{}, service.ErrUserNotFound
	}
	r.users[id] = user
	return user, nil
}
func (r *InMemoryUserRepo) Delete(id uuid.UUID) error {
	_, ok := r.users[id]
	if !ok {

		return service.ErrUserNotFound
	}
	delete(r.users, id)
	return nil

}

type InMemoryUserRepo struct {
	users map[uuid.UUID]service.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[uuid.UUID]service.User),
	}
}
