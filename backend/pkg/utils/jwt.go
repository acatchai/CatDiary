package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJwtSecret 获取JWT密钥
var getJwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "catdiary_default_secret_key_please_change"
	}
	return []byte(secret)
}

// GenerateToken 生成JWT
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat64, ok := claims["userID"].(float64)
		if !ok {
			return 0, errors.New("invalid user_id in token")
		}
		return uint(userIDFloat64), nil
	}
	return 0, errors.New("invalid token")
}
