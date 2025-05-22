package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"
	"time"
)

func (h *Handler) GetTask(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	task, err := h.repo.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) AddStudentsTask(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	input := models.StudentTask{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.TaskID = taskID
	input.StudentID = studentID.(int)

	studentTask, err := h.repo.AddStudentsTask(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go h.processTaskEvaluation(studentTask)
	c.JSON(http.StatusCreated, gin.H{"student_task": studentTask})
}

func (h *Handler) GetStudentTask(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	task, err := h.repo.GetStudentTask(taskID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) GetStudentTasks(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	tasks, err := h.repo.GetStudentTasks(taskID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *Handler) processTaskEvaluation(studentTask models.StudentTask) {
	// Создаем новый контекст для асинхронной обработки
	ctx := context.Background()

	task, err := h.repo.GetTask(studentTask.TaskID)
	if err != nil {
		log.Printf("Ошибка при получении задания: %v", err)
		return
	}
	studentCode := studentTask.Code

	taskDescription := task.Description

	evaluation, err := h.codeEvaluator.EvaluateCode(ctx, studentCode, taskDescription)
	if err != nil {
		log.Printf("Ошибка при оценке кода: %v", err)
		return
	}

	if evaluation.RequirementsScore == 0 {
		evaluation.BoundaryHandlingScore = 0
		evaluation.OptimizationScore = 0
		evaluation.ImplementationScore = 0
	}

	studentTask.Score = float32(float32(evaluation.RequirementsScore)*0.35 +
		float32(evaluation.BoundaryHandlingScore)*0.31 +
		float32(evaluation.OptimizationScore)*0.22 +
		float32(evaluation.ImplementationScore)*0.12)

	if evaluation.RequirementsScore == 5 && studentTask.Score >= 4.5 {
		studentTask.Status = "completed"
	} else {
		studentTask.Status = "canceled"
		studentTask.CodeWithCommentsByLLM = evaluation.CodeWithComment
	}

	studentTask.Requirements = evaluation.Requirements
	studentTask.RequirementsScore = evaluation.RequirementsScore

	studentTask.Implementation = evaluation.Implementation
	studentTask.ImplementationScore = evaluation.ImplementationScore

	studentTask.BoundaryHandling = evaluation.BoundaryHandling
	studentTask.BoundaryHandlingScore = evaluation.BoundaryHandlingScore

	studentTask.Optimization = evaluation.Optimization
	studentTask.OptimizationScore = evaluation.OptimizationScore

	studentTask.Recommendation = evaluation.Conclusion

	err = h.repo.UpdateStudentTask(studentTask)
	if err != nil {
		log.Printf("Ошибка при сохранении оценки: %v", err)
		return
	}
	log.Printf("Успешно оценено задание ID: %d, студент ID: %d, оценка: %.1f",
		studentTask.TaskID, studentTask.StudentID, studentTask.Score)

	log.Printf("New Code : %s ", evaluation.CodeWithComment)
	//log.Printf("Optimization : %s ", studentTask.Optimization)
	//log.Printf("Implementation: %s ", studentTask.Implementation)
	//log.Printf("BoundaryHandling: %s ", studentTask.BoundaryHandling)
	//log.Printf("recomends: %d ", evaluation.Conclusion)

}

func (h *Handler) GetFinalScore(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}

	score, err := h.repo.GetFinalScore(taskID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"score": score})
}

func (h *Handler) AddTask(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		log.Println("Unauthorized")
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		log.Println("Invalid module id")
		return
	}

	input := models.Task{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while binding json: " + err.Error())
		return
	}
	input.ModuleID = moduleID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	task, err := h.repo.AddTask(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while adding task: " + err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": task})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		log.Println("Invalid task id")
		return
	}

	err = h.repo.DeleteTask(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while deleting task: " + err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateTask(c *gin.Context) {
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

	if teacher.IsTeacher != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		log.Println("Invalid task id")
		return
	}

	task, err := h.repo.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while getting task: " + err.Error())
		return
	}

	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while binding json: " + err.Error())
	}

	if name, ok := input["name"].(string); ok {
		task.Name = name
	}
	if desc, ok := input["description"].(string); ok {
		task.Description = desc
	}

	if deadlineStr, ok := input["deadline"].(string); ok {
		deadline, err := time.Parse(time.RFC3339, deadlineStr)
		if err == nil {
			task.Deadline = deadline
		}
	}

	if is_active, ok := input["is_active"].(bool); ok {
		task.IsActive = is_active
	}

	err = h.repo.UpdateTask(task)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while updating task: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) GetAllTasks(c *gin.Context) {
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

	if teacher.IsTeacher != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *Handler) GetTaskAnswersForTeacher(c *gin.Context) {
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
	if teacher.IsTeacher != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		log.Println("Invalid task id")
		return
	}
	task, err := h.repo.GetTaskAnswersForTeacher(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error while getting task: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
