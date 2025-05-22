package handler

import (
	"fmt"
	"log"
	"net/http"
	"onlineschool/internal/auth"
	"onlineschool/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Регистрация пользователя
func (h *Handler) Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.repo.FindByEmail(input.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
		return
	}

	if err := h.repo.FindByPhone(input.Phone); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким номером телефона уже существует"})
		return
	}

	// Генерация соли
	salt, err := auth.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации соли"})
		return
	}

	log.Println("Хешируемый пароль и соль: ", input.Password, salt)
	// Хеширование пароля
	hashedPassword := auth.HashPassword(input.Password, salt)

	// Присваиваем соли и хешированный пароль пользователю
	// fmt.Println(input)
	input.Password = hashedPassword
	input.Salt = salt
	// fmt.Println(salt)

	// Сохраняем пользователя в базу данных
	err = h.repo.CreateUser(&input)
	// log.Println("ok")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении пользователя"})
		return
	}

	createdUser, err := h.repo.FindByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске пользователя"})
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(uint(createdUser.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован", "login": createdUser.Email,
		"userId": createdUser.ID,
		"token":  token})
}

// Логин пользователя
func (h *Handler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Привязываем данные из запроса
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя по email
	user, err := h.repo.FindByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	log.Println("Соль из бд:", user.Salt)
	log.Println("Введенный пароль", input.Password)
	// Проверяем пароль
	if !auth.CheckPassword(input.Password, user.Salt, user.Password) {
		fmt.Println("oou")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "клиент успешно авторизован",
		"login":      user.Email,
		"userId":     user.ID,
		"token":      token,
		"is_teacher": user.IsTeacher, // Возвращаем токен в ответе
	})
}

func (h *Handler) Logout(c *gin.Context) {

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	if len(token) > 7 && strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}

	claims, err := auth.ValidateJWT(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed: empty claims"})
		return
	}

	if err := claims.Valid(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or invalid"})
		return
	}

	expFloat, ok := (*claims)["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exp not found in token or has wrong type"})
		return
	}

	expTime := time.Unix(int64(expFloat), 0)

	tokenTTL := time.Until(expTime)

	if tokenTTL <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token has already expired"})
		return
	}

	err = h.redisClient.WriteJWTToBlacklist(c.Request.Context(), token, tokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выходе из системы"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Выход из системы успешно выполнен"})
}

// TokenAuth middleware для проверки JWT токена
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
			c.Abort()
			return
		}

		if len(token) > 7 && strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed: empty claims"})
			c.Abort()
			return
		}

		userIDRaw, ok := (*claims)["user_id"]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
			c.Abort()
			return
		}

		var userID int
		switch v := userIDRaw.(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		case string:
			var parseErr error
			userID, parseErr = strconv.Atoi(v)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
				c.Abort()
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected user ID type"})
			c.Abort()
			return
		}

		// Устанавливаем userId в контекст
		c.Set("userId", userID)
		log.Println("userId from token: ", userID)

		c.Next()
	}
}
