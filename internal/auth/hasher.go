package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateSalt создаёт случайную соль
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword хеширует пароль с солью
func HashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// CheckPassword проверяет пароль
func CheckPassword(password, salt, hash string) bool {
	fmt.Println(HashPassword(password, salt), hash)
	return HashPassword(password, salt) == hash
}
