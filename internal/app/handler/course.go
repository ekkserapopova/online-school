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

func (h *Handler) GetCourses(c *gin.Context) {
	name := c.Query("name")
	courses, err := h.repo.GetCourses(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (h *Handler) GetCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
	}
	course, err := h.repo.GetCourse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) GetStudentsCourses(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//user, _ := h.repo.GetByID(userID.(int))
	courses, err := h.repo.GetStudentsCourses(studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
	if len(courses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No courses found"})
		return
	}
}

func (h *Handler) GetTeachersCourses(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//user, _ := h.repo.GetByID(userID.(int))
	courses, err := h.repo.GetTeachersCourses(teacherID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
	if len(courses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No courses found"})
		return
	}
}

func (h *Handler) GetStudentsCourse(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.repo.GetStudentsCourse(studentID.(int), courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) GetTeachersCourse(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.repo.GetTeachersCourse(teacherID.(int), courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) GetLanguages(c *gin.Context) {
	name := c.Query("name")

	langs, err := h.repo.GetLanguages(name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"languages": langs})
}

func (h *Handler) EnrollStudent(c *gin.Context) {

	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//user, _ := h.repo.GetByID(studentID.(int))

	courseID, err := strconv.Atoi(c.Param("courseID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	err = h.repo.EnrollStudent(studentID.(int), courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student enrolled successfully"})

}

func (h *Handler) IsStudentEnrolledInCourse(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//user, _ := h.repo.GetByID(userID.(int))

	courseID, err := strconv.Atoi(c.Param("courseID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	isEnrolled, err := h.repo.IsStudentEnrolledInCourse(studentID.(int), courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_enrolled": isEnrolled})
}

func (h *Handler) GetUserProgressForAllCourses(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	progressCourses, err := h.repo.GetUserProgressForAllCourses(studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"progress": progressCourses})
}

func (h *Handler) CreateCourse(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	input := models.Course{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	input.TeacherID = teacherID.(int)
	input.IsActive = true

	course, err := h.repo.AddCourse(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (h *Handler) GetCoursePhoto(c *gin.Context) {
	//_, ok := c.Get("userId")
	//if !ok {
	//	log.Println("Not authorized")
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	course, err := h.repo.GetCourse(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	photoPath := course.Image
	log.Println("photoPath:", photoPath)
	if photoPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo not found"})
		return
	}
	log.Print("photo path: ", photoPath)

	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file does not exist"})
		return
	}

	// Отдаём файл с правильными заголовками
	c.File(photoPath)
}

func (h *Handler) UpdateCoursePhoto(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		log.Println("Not authorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.repo.GetCourse(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course not found"})
		return
	}

	// Получение файла из запроса
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get photo from request"})
		return
	}
	defer file.Close()

	// Создание пути для сохранения нового изображения
	uploadDir := "./photos/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	fileExt := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("course_%d_%s%s", courseID, uuid.New().String(), fileExt)
	newFilePath := filepath.Join(uploadDir, newFileName)

	// Сохранение файла
	out, err := os.Create(newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write photo to disk"})
		return
	}

	// Обновление информации о курсе (путь к изображению)
	course.Image = newFilePath
	if err := h.repo.UpdateCourseImage(course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update course image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course photo updated successfully"})
}
