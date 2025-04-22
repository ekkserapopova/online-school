package handler

import (
	"net/http"
	"strconv"

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

	lessonID, err := strconv.Atoi(c.Param("id"))
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
