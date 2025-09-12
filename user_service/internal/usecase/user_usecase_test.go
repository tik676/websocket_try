package usecase

import (
	"errors"
	"testing"
	"time"
	"user_service/internal/domain"
	"user_service/internal/domain/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

func TestUseCase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthorization(ctrl)
	mockToken := mocks.NewMockTokenManager(ctrl)
	usecase := NewUseCase(mockAuth, mockToken)

	t.Run("success", func(t *testing.T) {
		password := "password123"
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: password,
		}

		userFromdb := &domain.User{
			ID:           1,
			Name:         "testuser",
			Role:         "user",
			PasswordHash: string(hashPassword),
		}

		expected := &domain.Token{
			AccessToken:  "acces_token",
			RefreshToken: "refresh_token",
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		}

		mockAuth.EXPECT().
			Login(gomock.Eq(input)).
			Return(userFromdb, nil).
			Times(1)

		mockToken.EXPECT().
			CreateToken(userFromdb.ID, userFromdb.Name, userFromdb.Role).
			Return(expected, nil).
			Times(1)

		token, err := usecase.LoginUser(input)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, expected.AccessToken, token.AccessToken)
		assert.Equal(t, expected.RefreshToken, token.RefreshToken)
	})

	t.Run("empty_name", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "",
			Password: "password123",
		}

		token, err := usecase.LoginUser(input)

		assert.Error(t, err)
		assert.Nil(t, token)

	})

	t.Run("empty_password", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: "",
		}

		token, err := usecase.LoginUser(input)

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("user_not_found", func(t *testing.T) {
		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: "password123",
		}

		mockAuth.EXPECT().
			Login(gomock.Eq(input)).
			Return(nil, errors.New("user not found")).
			Times(1)

		token, err := usecase.LoginUser(input)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("invalid_password", func(t *testing.T) {
		wrondPassword := "wrong_pass"
		correctPassword := "correct_pass"

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: wrondPassword,
		}

		userFromdb := &domain.User{
			ID:           1,
			Name:         "testuser",
			Role:         "user",
			PasswordHash: string(hashPassword),
		}

		mockAuth.EXPECT().
			Login(gomock.Eq(input)).
			Return(userFromdb, nil).
			Times(1)

		token, err := usecase.LoginUser(input)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, "invalid password", err.Error())
	})

	t.Run("token_creation_failed", func(t *testing.T) {
		password := "password123"
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		input := domain.AuthorizationInput{
			Name:     "testuser",
			Password: password,
		}

		userFromdb := &domain.User{
			ID:           1,
			Name:         "testuser",
			Role:         "user",
			PasswordHash: string(hashPassword),
		}

		mockAuth.EXPECT().
			Login(gomock.Eq(input)).
			Return(userFromdb, nil).
			Times(1)

		mockToken.EXPECT().
			CreateToken(userFromdb.ID, userFromdb.Name, userFromdb.Role).
			Return(nil, errors.New("token service unavailable")).
			Times(1)

		token, err := usecase.LoginUser(input)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, "failed to create token", err.Error())
	})
}

func TestUseCase_LoginAnonUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthorization(ctrl)
	mockToken := mocks.NewMockTokenManager(ctrl)
	usecase := NewUseCase(mockAuth, mockToken)

	t.Run("success", func(t *testing.T) {
		expectedToken := &domain.Token{
			AccessToken:  "success_token",
			RefreshToken: "success_refresh_token",
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		}

		mockToken.EXPECT().
			CreateToken(gomock.Any(), gomock.Any(), gomock.Eq("anon")).
			Return(expectedToken, nil).
			Times(1)

		token, err := usecase.LoginAnonUser()

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token, expectedToken)
	})

	t.Run("failed_create_token", func(t *testing.T) {
		mockToken.EXPECT().
			CreateToken(gomock.Any(), gomock.Any(), gomock.Eq("anon")).
			Return(nil, errors.New("Failed work token manager")).
			Times(1)

		token, err := usecase.LoginAnonUser()

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, "Failed to create token", err.Error())
	})
}
