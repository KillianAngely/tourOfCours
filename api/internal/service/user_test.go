package service

import (
	"api/internal/domain"
	"api/internal/service/mock"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestUserService(t *testing.T) {

	t.Run("Create a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockUserRepository(ctrl)
		s := NewUserService(repo)

		input := CreateUserInput{
			FirstName: "john",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
		}

		expected := domain.User{
			ID:        uuid.New(),
			FirstName: "john",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
		}

		repo.EXPECT().Create(gomock.Cond(func(u domain.User) bool {
			return u.FirstName == input.FirstName &&
				u.LastName == input.LastName &&
				u.Email == input.Email
		})).Return(expected, nil)

		_, err := s.Create(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

}
