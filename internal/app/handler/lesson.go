package handler

import (
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLessons(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	period := c.Query("period")

	if period == "" {
		period = "all"
	}

	validPeriods := map[string]bool{
		"past":   true,
		"future": true,
		"all":    true,
	}

	if _, ok := validPeriods[period]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period parameter"})
		return
	}

	lessons, err := h.repo.GetLessons(studentID.(int), period)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}

func (h *Handler) GetLesson(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
	}
	lesson, err := h.repo.GetLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"lesson": lesson})

}

func (h *Handler) GetLessonsByCourseID(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	lessons, err := h.repo.GetLessonsByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}

func (h *Handler) AddLesson(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		log.Println("unauthorized")
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		log.Println("Invalid module id")
		return
	}
	input := models.Lesson{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Invalid JSON")
		return
	}

	input.ModuleID = moduleID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	lesson, err := h.repo.AddLesson(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Failed to add lesson")
		return
	}
	c.JSON(http.StatusOK, gin.H{"lesson": lesson})
}

func (h *Handler) DeleteLesson(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	lessonID, err := strconv.Atoi(c.Param("lessonID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
		return
	}

	err = h.repo.DeleteLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	lessonID, err := strconv.Atoi(c.Param("lessonID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
		return
	}

	lesson, err := h.repo.GetLesson(lessonID)

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if name, ok := input["name"].(string); ok {
		lesson.Name = name
	}
	if desc, ok := input["description"].(string); ok {
		lesson.Description = desc
	}
	if startStr, ok := input["start"].(string); ok {
		start, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			lesson.Start = start
		}
	}
	if endStr, ok := input["end"].(string); ok {
		end, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			lesson.End = end
		}
	}

	if is_active, ok := input["is_active"].(bool); ok {
		lesson.IsActive = is_active
	}
	lesson.UpdatedAt = time.Now()

	updatedLesson, err := h.repo.UpdateLesson(lesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lesson": updatedLesson})
}
