package utils

import (
	"errors"
	"os"
	"time"

	models "github.com/Rifq11/Trava-be/models"
	"github.com/golang-jwt/jwt/v5"
)

var defaultJWTSecret = []byte("trava-secret")

func getJWTSecret() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	return defaultJWTSecret
}

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   int    `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
