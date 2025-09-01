package usecase

import (
	"errors"
	"fmt"
	"user_service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo      domain.Authorization
	repoToken domain.TokenManager
}

func NewUseCase(repo domain.Authorization, token domain.TokenManager) *UseCase {
	return &UseCase{repo: repo, repoToken: token}
}

func (u *UseCase) RegisterUser(name, password string) (*domain.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.Register(domain.AuthorizationInput{
		Name:     name,
		Password: string(hashPassword),
		Role:     "user",
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UseCase) LoginUser(name, password string) (*domain.Token, error) {
	user, err := u.repo.Login(domain.AuthorizationInput{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := u.repoToken.CreateToken(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	return token, nil
}

func (u *UseCase) LoginAnonUser() (*domain.Token, error) {
	anonName := fmt.Sprintf("anon_%s", uuid.New().String())

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)

	user, err := u.repo.Register(domain.AuthorizationInput{
		Name:     anonName,
		Password: string(hashPassword),
		Role:     "anon",
	})
	if err != nil {
		return nil, err
	}

	token, err := u.repoToken.CreateToken(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	return token, nil
}

func (u *UseCase) RefreshToken(refreshToken string) (*domain.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	newTokens, err := u.repoToken.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}

	return newTokens, nil
}

func (u *UseCase) LogoutUser(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}

	err := u.repoToken.RevokeRefreshToken(refreshToken)
	if err != nil {
		return errors.New("refresh token not found")
	}

	return nil
}
