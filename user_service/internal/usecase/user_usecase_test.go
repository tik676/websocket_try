package usecase

import (
	"errors"
	"testing"
	"user_service/internal/domain"
	"user_service/internal/domain/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthorization(ctrl)
	mockToken := mocks.NewMockTokenManager(ctrl)
	usecase := NewUseCase(mockAuth, mockToken)

	t.Run("succes", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: "password123",
		}

		expectedUser := &domain.User{
			ID:   1,
			Name: "testuser",
			Role: "user",
		}

		mockAuth.EXPECT().
			Register(gomock.Any()).
			Return(expectedUser, nil).
			Times(1)

		user, err := usecase.RegisterUser(input)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, input.Name, user.Name)
	})

	t.Run("bcrypt_correctly", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "correctPass_user",
			Password: "password123",
		}

		var capturedInput domain.AuthorizationInput

		mockAuth.EXPECT().
			Register(gomock.Any()).
			DoAndReturn(func(regInput domain.AuthorizationInput) (*domain.User, error) {
				capturedInput = regInput
				return &domain.User{ID: 1, Name: regInput.Name, Role: "user"}, nil
			}).
			Times(1)

		_, err := usecase.RegisterUser(input)

		assert.NoError(t, err)

		assert.NotEqual(t, input.Password, capturedInput.Password)
		assert.True(t, len(capturedInput.Password) > 20)

		assert.Equal(t, input.Name, capturedInput.Name)

	})

	t.Run("bcrypt_error", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "badPass_user",
			Password: "mqvdrradwxmbstwrvbzjkuvqendprenemxfkvydzrymofxjnbwhjtnfkusvkywtzsmrvyrbehzfrthtt",
		}

		user, err := usecase.RegisterUser(input)

		assert.Error(t, err)
		assert.Nil(t, user)

	})

	t.Run("repository_error", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: "password123",
		}

		mockAuth.EXPECT().
			Register(gomock.Any()).
			Return(nil, errors.New("user already exists")).
			Times(1)

		user, err := usecase.RegisterUser(input)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user already exists")
	})

	t.Run("name_empty", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "",
			Password: "password123",
		}

		user, err := usecase.RegisterUser(input)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("password_empty", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: "",
		}

		user, err := usecase.RegisterUser(input)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
