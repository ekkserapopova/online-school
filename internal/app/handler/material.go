package handler

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"onlineschool/internal/models"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMaterials(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
	}
	materials, err := h.repo.GetMaterials(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"materials": materials})
}

func (h *Handler) GetMaterial(c *gin.Context) {
	materialID, err := strconv.Atoi(c.Param("materialID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material id"})
	}
	material, err := h.repo.GetMaterial(materialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"material": material})
}

func (h *Handler) GetMaterialPDF(c *gin.Context) {
	//_, ok := c.Get("userId")
	//if !ok {
	//	log.Println("Not authorized")
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}

	materialID, err := strconv.Atoi(c.Param("materialID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material id"})
		return
	}

	material, err := h.repo.GetMaterial(materialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfPath := filepath.Join("lessonMaterials", material.File) // если у тебя папка "files"
	log.Println("pdfPath:", pdfPath)

	if pdfPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found"})
		return
	}

	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF file does not exist"})
		return
	}

	// Устанавливаем правильный Content-Type для PDF
	c.Header("Content-Type", "application/pdf")
	c.File(pdfPath)
}

func (h *Handler) CreateMaterial(c *gin.Context) {
	log.Println("CreateMaterial: start")

	teacherID, ok := c.Get("userId")
	if !ok {
		log.Println("CreateMaterial: unauthorized - userId not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("CreateMaterial: teacherID = %v\n", teacherID)

	teacher, err := h.repo.GetByID(teacherID.(int))
	if err != nil {
		log.Printf("CreateMaterial: failed to get teacher by ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !teacher.IsTeacher {
		log.Println("CreateMaterial: forbidden - user is not a teacher")
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		log.Printf("CreateMaterial: invalid lesson ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("CreateMaterial: lessonID = %d\n", lessonID)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("CreateMaterial: failed to retrieve file from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()
	log.Printf("CreateMaterial: received file: %s\n", header.Filename)

	uploadDir := "./lessonMaterials/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("CreateMaterial: failed to create directory: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create upload directory"})
		return
	}

	fileExt := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("material_%d_%s%s", teacherID, uuid.New().String(), fileExt)
	newFilePath := filepath.Join(uploadDir, newFileName)

	out, err := os.Create(newFilePath)
	if err != nil {
		log.Printf("CreateMaterial: failed to create file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		log.Printf("CreateMaterial: failed to write file to disk: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file to disk"})
		return
	}
	log.Printf("CreateMaterial: file saved to %s\n", newFilePath)

	var input models.Material
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("CreateMaterial: failed to bind input: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.File = newFilePath
	input.LessonID = lessonID
	input.Name = c.PostForm("name")

	err = h.repo.CreateMaterial(input)
	if err != nil {
		log.Printf("CreateMaterial: failed to create material in repo: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CreateMaterial: material created successfully: %+v\n", input)
	c.JSON(http.StatusOK, gin.H{"material": input})
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teacher, err := h.repo.GetByID(teacherID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !teacher.IsTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	_, err = strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		log.Printf("CreateMaterial: invalid lesson ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	materialID, err := strconv.Atoi(c.Param("materialID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material id"})
		return
	}

	err = h.repo.DeleteMaterial(materialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) UpdateMaterial(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teacher, err := h.repo.GetByID(teacherID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !teacher.IsTeacher {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	_, err = strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		log.Printf("CreateMaterial: invalid lesson ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	materialID, err := strconv.Atoi(c.Param("materialID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material id"})
		return
	}

	material, err := h.repo.GetMaterial(materialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if name, ok := input["name"].(string); ok {
		material.Name = name
	}

	err = h.repo.UpdateMaterial(material)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"material": material})
}
