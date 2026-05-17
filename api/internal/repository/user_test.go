package repository

import (
	"api/internal/domain"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestUser() domain.User {
	return domain.User{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
	}
}

func setupRepoWithUser(t *testing.T) (*InMemoryUserRepo, domain.User) {
	t.Helper()

	r := NewInMemoryUserRepo()
	u := newTestUser()

	if _, err := r.Create(u); err != nil {
		t.Fatalf("Setup: Create returned an unexpected error: %v", err)
	}

	return r, u
}

func TestInMemoryUserRepo(t *testing.T) {
	t.Run("creates and finds a user by ID", func(t *testing.T) {
		r := NewInMemoryUserRepo()

		u := newTestUser()

		_, errCreate := r.Create(u)
		if errCreate != nil {
			t.Fatalf("Create returned an unexpected error: %v", errCreate)
		}

		foundUser, errFind := r.FindByID(u.ID)
		if errFind != nil {
			t.Fatalf("FindByID returned an unexpected error: %v", errFind)
		}
		if foundUser != u {
			t.Errorf("FindByID returned wrong user: got %+v, want %+v", foundUser, u)
		}
	})

	t.Run("Updates an existing user", func(t *testing.T) {

		r, original := setupRepoWithUser(t)

		updated := original
		updated.FirstName = "josef"

		returnedUser, err := r.Update(updated.ID, updated)
		if err != nil {
			t.Fatalf("Update returned an unexpected error: %v", err)
		}

		if returnedUser.FirstName != "josef" {
			t.Errorf("Update did not change FirstName: got %q, want %q", returnedUser.FirstName, "josef")
		}

		if returnedUser.ID != original.ID {
			t.Errorf("Update changed ID: got %v, want %v", returnedUser.ID, original.ID)
		}
		if returnedUser.LastName != original.LastName {
			t.Errorf("Update changed LastName: got %q, want %q", returnedUser.LastName, original.LastName)
		}
		if !returnedUser.CreatedAt.Equal(original.CreatedAt) {
			t.Errorf("Update changed CreatedAt: got %v, want %v", returnedUser.CreatedAt, original.CreatedAt)
		}

		persistedUser, err := r.FindByID(updated.ID)
		if err != nil {
			t.Fatalf("FindByID returned an unexpected error after Update: %v", err)
		}
		if !reflect.DeepEqual(persistedUser, updated) {
			t.Errorf("Update did not persist correctly: got %+v, want %+v", persistedUser, updated)
		}
	})

	t.Run("deletes a user", func(t *testing.T) {
		r, u := setupRepoWithUser(t)

		errDelete := r.Delete(u.ID)
		if errDelete != nil {
			t.Fatalf("Delete returned an unexpected error: %v", errDelete)

		}

		_, errFind := r.FindByID(u.ID)
		if !errors.Is(errFind, ErrUserNotFound) {
			t.Fatalf("FindById after Delete: returned an unexpected error: %v", errFind)
		}

	})

}

func TestInMemoryUserRepo_NotFoundErrors(t *testing.T) {

	cases := []struct {
		name string
		op   func(*InMemoryUserRepo, uuid.UUID) error
	}{
		{
			name: "FindByID returns ErrUserNotFound",
			op: func(r *InMemoryUserRepo, id uuid.UUID) error {
				_, err := r.FindByID(id)
				return err
			},
		},
		{
			name: "Delete returns ErrUserNotFound",
			op: func(r *InMemoryUserRepo, id uuid.UUID) error {
				return r.Delete(id)
			},
		},
		{
			name: "Update returns ErrUserNotFound",
			op: func(r *InMemoryUserRepo, id uuid.UUID) error {
				_, err := r.Update(id, domain.User{})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewInMemoryUserRepo()
			err := tc.op(r, uuid.New())

			if !errors.Is(err, ErrUserNotFound) {
				t.Errorf("Expected ErrUserNotFound, got %v", err)
			}
		})
	}
}
