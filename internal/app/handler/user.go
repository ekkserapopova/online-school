package handler

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	user, err := h.repo.GetByID(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "name": user.Name, "surname": user.Surname, "email": user.Email, "is_teacher": user.IsTeacher})
}

func (h *Handler) AddPhoto(c *gin.Context) {
	// Ограничение размера файла (например, 50 МБ)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50<<20)

	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.repo.GetByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
		return
	}
	defer file.Close()

	const uploadDir = "./photos/users"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	//filename := filepath.Base(header.Filename) // защищаемся от ".."
	uniqueName := uuid.New().String() + filepath.Ext(header.Filename)
	destPath := filepath.Join(uploadDir, uniqueName)

	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	if err := h.repo.AddPhoto(user, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user photo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("файл %s успешно загружен", uniqueName)})
}

func (h *Handler) GetPhoto(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.repo.GetByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	photoPath := user.Photo
	if photoPath == "" {

		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	log.Print("photo path: ", photoPath)

	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo file does not exist"})
		return
	}

	// Отдаём файл с правильными заголовками
	c.File(photoPath)
}
