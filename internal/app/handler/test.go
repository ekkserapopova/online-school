package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
