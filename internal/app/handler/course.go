package handler

import (
	"log"
	"net/http"
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

func (h *Handler) GetUsersCourses(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}
	user, _ := h.repo.GetByID(userID.(int))
	courses, err := h.repo.GetUsersCourses(user.ID)
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

func (h *Handler) GetCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
	}
	course, err := h.repo.GetCourse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
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

	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	user, _ := h.repo.GetByID(userID.(int))

	courseID, err := strconv.Atoi(c.Param("courseID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	err = h.repo.EnrollStudent(user.ID, courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student enrolled successfully"})

}

func (h *Handler) IsStudentEnrolledInCourse(c *gin.Context) {
	userID, ok := c.Get("userId")
	log.Println(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	user, _ := h.repo.GetByID(userID.(int))

	courseID, err := strconv.Atoi(c.Param("courseID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	isEnrolled, err := h.repo.IsStudentEnrolledInCourse(user.ID, courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_enrolled": isEnrolled})
}
