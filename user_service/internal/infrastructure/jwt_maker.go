package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"user_service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
	DB        *sql.DB
}

func NewJWTMaker(secretKey string, db *sql.DB) *JWTMaker {
	return &JWTMaker{secretKey: secretKey, DB: db}
}

func (j *JWTMaker) CreateToken(userID int64, name, role string) (*domain.Token, error) {
	now := time.Now()
	access_expires_at := now.Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"name":    name,
		"exp":     access_expires_at.Unix(),
		"iat":     now.Unix(),
		"iss":     "user-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accesstoken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Printf("Failed to create token:%v", err)
		return nil, err
	}

	refreshToken := j.GenerateRefreshToken()
	expires_at := now.Add(7 * 24 * time.Hour)
	query := `INSERT INTO tokens(refresh_token,created_at,expires_at,user_id)VALUES($1,$2,$3,$4)`
	_, err = j.DB.Exec(query, refreshToken, now, expires_at, userID)
	if err != nil {
		log.Printf("Failed to save token:%v", err)
		return nil, err
	}

	return &domain.Token{
		AccessToken:  accesstoken,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		ExpiresAt:    expires_at,
	}, nil

}

func (j *JWTMaker) VerifyToken(tokenString string) (userID int64, name, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = int64(claims["user_id"].(float64))
		role = claims["role"].(string)
		name = claims["name"].(string)
		return userID, name, role, nil
	}

	return 0, "", "", errors.New("invalid token claims")
}

func (j *JWTMaker) RevokeRefreshToken(refreshToken string) error {
	query := `DELETE FROM tokens WHERE refresh_token=$1`
	_, err := j.DB.Exec(query, refreshToken)
	if err != nil {
		log.Printf("error token not found:%v", err)
		return err
	}
	return nil
}

func (j *JWTMaker) GenerateRefreshToken() string {
	refreshToken := uuid.New().String()
	return refreshToken

}

func (j *JWTMaker) RefreshAccessToken(refreshToken string) (*domain.Token, error) {

	userID, err := j.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var role string

	userQuery := `SELECT role FROM users WHERE id = $1`
	err = j.DB.QueryRow(userQuery, userID).Scan(&role)
	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, errors.New("user not found")
	}

	now := time.Now()
	accessExpiresAt := now.Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     accessExpiresAt.Unix(),
		"iat":     now.Unix(),
		"iss":     "user-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	newAccessToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Printf("Failed to create new access token: %v", err)
		return nil, err
	}
	return &domain.Token{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		ExpiresAt:    accessExpiresAt,
	}, nil
}

func (j *JWTMaker) VerifyRefreshToken(refreshToken string) (int64, error) {
	var userID int64
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM tokens WHERE refresh_token = $1`
	err := j.DB.QueryRow(query, refreshToken).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, errors.New("refresh token not found")
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("refresh token expired")
	}

	return userID, nil
}
