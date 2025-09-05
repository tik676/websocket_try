package usecase

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"user_service/internal/domain"

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
	anonID := rand.Int64()
	anonName := fmt.Sprintf("anon_%v", anonID)

	user := &domain.User{
		ID:   anonID,
		Name: anonName,
		Role: "anon",
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
