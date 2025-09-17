package domain

import (
	"context"
	"time"
)

//go:generate mockgen -source=user.go -destination=mocks/user.go -package=mocksauth

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	RegisteredAt time.Time `json:"registered_at"`
}

type AuthorizationInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserActionEvent struct {
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Timestamp time.Time `json:"timestamp"`
}

type Authorization interface {
	Register(input AuthorizationInput) (*User, error)
	Login(input AuthorizationInput) (*User, error)
}

type TokenManager interface {
	CreateToken(userID int64, name, role string) (*Token, error)
	VerifyToken(token string) (userID int64, name, role string, err error)

	GenerateRefreshToken() string
	RefreshAccessToken(refreshToken string) (*Token, error)
	VerifyRefreshToken(refreshToken string) (int64, error)
	RevokeRefreshToken(refreshToken string) error
}

type KafkaProducer interface {
	SendUserRegistered(ctx context.Context, userID int64, name string) error
	SendUserLoggedIn(ctx context.Context, userID int64, name string) error
	Close() error
}
