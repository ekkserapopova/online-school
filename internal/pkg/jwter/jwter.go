package jwter

import (
	"fmt"
	"go.uber.org/fx"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/golang-jwt/jwt/v5"
)

type Params struct {
	fx.In

	Logger *slog.Logger
}

type JWTer struct {
	log *slog.Logger
}

func New(p Params) *JWTer {
	return &JWTer{
		log: p.Logger,
	}
}

// Секретный ключ для подписи токена (должен быть защищен)
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Функция для генерации JWT
func (jwter *JWTer) GenerateJWT(userID int) (string, error) {
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

// ValidateJWT проверяет токен и возвращает claims
func (jwter *JWTer) ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи: %v", token.Method)
		}
		return secretKey, nil
	})

	if err != nil {
		jwter.log.Error("Failed to parse token", "error", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("недействительный токен")
}

// GetUserID извлекает ID пользователя из JWT claims
func (jwter *JWTer) GetUserID(claims *jwt.MapClaims) (int, error) {
	userID, ok := (*claims)["user_id"]
	if !ok {
		return 0, fmt.Errorf("в токене отсутствует поле user_id")
	}

	fmt.Println("COMPLETED")

	// Преобразование userID в uint
	switch v := userID.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("некорректный тип для user_id: %T", userID)
	}
}
