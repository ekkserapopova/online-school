package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи токена (должен быть защищен)
var secretKey = []byte("katyushka")

// Функция для генерации JWT
func GenerateJWT(userID uint) (string, error) {
	// Создание claims (данных внутри токена)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Истечение через 24 часа
	}

	// Создание токена с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписание токена с использованием секретного ключа
	return token.SignedString(secretKey)
}

// ValidateJWT проверяет токен
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("недействительный токен")
}
