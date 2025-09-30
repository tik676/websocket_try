package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
	"user_service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo          domain.Authorization
	repoToken     domain.TokenManager
	eventProducer domain.KafkaProducer
}

func NewUseCase(repo domain.Authorization, token domain.TokenManager, repoEventProducer domain.KafkaProducer) *UseCase {
	return &UseCase{repo: repo,
		repoToken:     token,
		eventProducer: repoEventProducer,
	}
}

func (u *UseCase) RegisterUser(input domain.AuthorizationInput) (*domain.User, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.Password == "" {
		return nil, errors.New("password is required")
	}
	if input.Role == "" {
		input.Role = "user"
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashPassword)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.Register(input)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := u.eventProducer.SendUserRegistered(ctx, user.ID, user.Name); err != nil {
		log.Printf("Failed to send Kafka event: %v", err)
	}

	return user, nil
}

func (u *UseCase) LoginUser(input domain.AuthorizationInput) (*domain.Token, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.Password == "" {
		return nil, errors.New("password is required")
	}
	user, err := u.repo.Login(input)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := u.repoToken.CreateToken(user.ID, user.Name, user.Role)
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := u.eventProducer.SendUserLoggedIn(ctx, user.ID, user.Name); err != nil {
		log.Printf("Failed to send Kafka event: %v", err)
	}

	return token, nil
}

func (u *UseCase) LoginAnonUser() (*domain.Token, error) {
	anonID := rand.Int64()
	anonName := fmt.Sprintf("anon_%v", anonID)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(anonName), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := u.repo.Register(domain.AuthorizationInput{
		Name:     anonName,
		Password: string(hashedPassword),
		Role:     "anon",
	})

	token, err := u.repoToken.CreateToken(user.ID, user.Name, "anon")
	if err != nil {
		return nil, errors.New("failed to create token")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := u.eventProducer.SendUserLoggedIn(ctx, user.ID, user.Name); err != nil {
		log.Printf("Failed to send Kafka event: %v", err)
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
