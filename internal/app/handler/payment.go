package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPayment(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	payment, err := h.repo.GetPayment(userID.(int), courseID)
	if err != nil {
		log.Println("Error getting payment")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

func (h *Handler) AddPayment(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		log.Fatal("User ID not found in token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}
	log.Println("User ID found in token")

	user, _ := h.repo.GetByID(userID.(int))

	courseID, err := strconv.Atoi(c.Param("courseID"))

	if err != nil {
		log.Fatal("Invalid course id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}

	log.Println("Course ID found successful")

	payment, err := h.repo.AddPayment(user.ID, courseID)
	if err != nil {
		log.Fatal("Error in repo method 'AddPayment'")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	log.Println("Payment added successfully")
	c.JSON(http.StatusOK, gin.H{"payment": payment, "message": "Payment added successfully"})
}

func (h *Handler) UpdatePayment(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in token"})
		return
	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course id"})
		return
	}
	payment, err := h.repo.GetPayment(userID.(int), courseID)
	if err != nil {
		log.Println("Error getting payment")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// course, err := h.repo.GetCourse(courseID)
	// if err != nil {
	// 	log.Println("Error getting course")
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err})
	// 	return
	// }

	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	payment.Status = input.Status

	if err := h.repo.UpdatePaymentsStatus(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully"})

}
