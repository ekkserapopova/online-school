package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"
)

func (h *Handler) AddModuleToCourse(c *gin.Context) {
	//TODO: роль преподавателя
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return

	}

	courseID, err := strconv.Atoi(c.Param("courseID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	input := models.Module{}

	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	input.CourseID = courseID

	err = h.repo.AddModule(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"module": input})
}

func (h *Handler) GetModule(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	module, err := h.repo.GetModule(moduleID)

	log.Println(module)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"module": module})
}

func (h *Handler) DeleteModule(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teacher, err := h.repo.GetByID(teacherID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if teacher.IsTeacher != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "not teacher"})
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.repo.DeleteModule(moduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"module": nil})
}

func (h *Handler) UpdateModule(c *gin.Context) {
	teacherID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		log.Println("Unauthorized user")
		return
	}
	teacher, err := h.repo.GetByID(teacherID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Not user ID")
		return
	}

	if teacher.IsTeacher != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "not teacher"})
		log.Println("Not teacher")
		return
	}
	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Not module ID")
		return
	}

	module, err := h.repo.GetModule(moduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Not module ID")
		return
	}

	log.Println(module)

	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Not module input")
		return
	}
	log.Println(input)
	if name, ok := input["name"].(string); ok {
		log.Println(name)
		module.Name = name
	}

	if description, ok := input["description"].(string); ok {
		module.Description = description
	}

	log.Println(module.Name)
	err = h.repo.UpdateModule(module)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("Error updating module")
		return
	}

	c.JSON(http.StatusOK, gin.H{"module": module})
}
