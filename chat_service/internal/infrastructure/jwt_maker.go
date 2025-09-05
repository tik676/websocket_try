package infrastructure

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTmaker struct {
	secretKey string
}

func NewJWTmaker(key string) *JWTmaker {
	return &JWTmaker{secretKey: key}
}

func (jm *JWTmaker) VerifyToken(tokenString string) (userID int64, name, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return 0, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = int64(claims["user_id"].(float64))
		role = claims["role"].(string)
		name = claims["name"].(string)
		return userID, role, name, nil
	}
	return 0, "", "", errors.New("invalid token claims")
}
