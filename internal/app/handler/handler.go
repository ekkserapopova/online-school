package handler

import (
	"net/http"
	"onlineschool/internal/app"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo app.Repo
}

func NewHandler(repo app.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/teachers", h.GetTeachers)
	r.GET("/schedule/:id", h.GetSchedule)
	r.GET("/lesson/:id", h.GetLesson)
	r.GET("/courses", h.GetCourses)
	r.GET("/course/:id", h.GetCourse)
	r.GET("/student/:id", h.GetStudent)
	r.GET("/lesson/:id/materials", h.GetMaterials)
	r.GET("/course/:id/tests", h.GetTests)
	r.GET("/course/:id/test/:testid", h.GetTest)
	return r
}

func (h *Handler) GetTeachers(c *gin.Context) {
	name := c.Query("name")
	teachers, err := h.repo.GetTeachersList(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"teachers": teachers})
}

func (h *Handler) GetCourses(c *gin.Context) {
	name := c.Query("name")
	courses, err := h.repo.GetCourses(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
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

func (h *Handler) GetMaterials(c *gin.Context) {
	lessonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
	}
	materials, err := h.repo.GetMaterials(lessonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"materials": materials})

}

// func (h *Handler) GetLessons(c *gin.Context) {
// 	lessons, err := h.repo.GetLessons()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 	}
// 	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
// }

func (h *Handler) GetSchedule(c *gin.Context) {
	studentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
	}
	month, err := strconv.Atoi(c.Query("month"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
	}

	year, err := strconv.Atoi(c.Query("year"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	to := time.Date(year, time.Month(month)+1, 0, 23, 59, 59, 0, time.UTC)

	if to.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date"})
	}

	lessons, err := h.repo.GetSchedule(studentId, from, to)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"lessons": lessons})

}

func (h *Handler) GetStudent(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson id"})
	}
	user, err := h.repo.GetStudent(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

func (h *Handler) GetLesson(c *gin.Context) {
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

func (h *Handler) GetTest(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("testid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
	}

	test, err := h.repo.GetTest(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"test": test})
}

func (h *Handler) GetTests(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
	}

	tests, err := h.repo.GetTests(courseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"tests": tests})
}
