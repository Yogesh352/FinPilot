package util

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"stock-api/internal/config"
)

// var SecretKey = []byte(config.Load().JWTSecret)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(username string, userId int) (string, error) {
	secret := []byte(config.Load().JWTSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id": userId,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return jwt.NewWithClaims(jwt.SigningMethodHS256, token.Claims).SignedString(secret)
}

func ValidateToken(tokenStr string) (*jwt.Token, error) {
	secret := []byte(config.Load().JWTSecret)
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
}
