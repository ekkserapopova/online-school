package handler

import (
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTest(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	test, err := h.repo.GetTest(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"test": test})
}

func (h *Handler) GetTests(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
	}

	tests, err := h.repo.GetTests(courseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tests": tests})
}

func (h *Handler) GetQuestionsForTest(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	questions, err := h.repo.GetQuestionsForTest(testID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *Handler) AddAnswerByStudent(c *gin.Context) {

	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question id"})
		return
	}

	input := models.StudentAnswer{}

	input.QuestionID = questionID
	input.StudentID = userID.(int)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Student answer added: %+v", input)

	studentAnswer, err := h.repo.AddAnswerByStudent(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"answer": studentAnswer})
}

func (h *Handler) GetPointsOfTest(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	points, err := h.repo.GetPointsOfTest(testID, userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
}

func (h *Handler) GetRightAnswers(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	rightAnswers, err := h.repo.GetRightAnswers(testID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"right_answers": rightAnswers})
}

func (h *Handler) FinishTest(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	err = h.repo.FinishTest(testID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"test": testID})
}

func (h *Handler) GetStudentAnswers(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	answers, err := h.repo.GetStudentAnswers(testID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answers": answers})
}

func (h *Handler) GetCompletedTest(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
		return
	}

	completedTest, err := h.repo.GetCompletedTest(testID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"completed_test": completedTest})
}
