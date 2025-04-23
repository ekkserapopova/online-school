package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"
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
	c.JSON(http.StatusOK, gin.H{"student_task": studentTask})
}

func (h *Handler) processTaskEvaluation(studentTask models.StudentTask) {
	// Создаем новый контекст для асинхронной обработки
	ctx := context.Background()

	// Получаем информацию о задании, чтобы взять описание задачи
	task, err := h.repo.GetTask(studentTask.TaskID)
	if err != nil {
		log.Printf("Ошибка при получении задания: %v", err)
		return
	}
	// Получаем текст решения студента
	studentCode := studentTask.Code

	// Получаем описание задания
	taskDescription := task.Description

	// Вызываем оценку кода
	evaluation, err := h.codeEvaluator.EvaluateCode(ctx, studentCode, taskDescription)
	if err != nil {
		log.Printf("Ошибка при оценке кода: %v", err)
		return
	}

	studentTask.Score = evaluation.Score
	studentTask.Requirements = evaluation.Requirements
	studentTask.Implementation = evaluation.Implementation
	studentTask.BoundaryHandling = evaluation.BoundaryHandling
	studentTask.Optimization = evaluation.Optimization

	// Сохраняем результаты в БД
	err = h.repo.UpdateStudentTask(studentTask)
	if err != nil {
		log.Printf("Ошибка при сохранении оценки: %v", err)
		return
	}

	// Можно также отправить уведомление студенту о готовности оценки
	// например, через WebSocket или внутреннюю систему уведомлений
	log.Printf("Успешно оценено задание ID: %d, студент ID: %d, оценка: %.1f",
		studentTask.TaskID, studentTask.StudentID, evaluation.Score)

	log.Printf("Requirements : %s ", studentTask.Requirements)
	log.Printf("Optimization : %s ", studentTask.Optimization)
	log.Printf("Implementation: %s ", studentTask.Implementation)
	log.Printf("BoundaryHandling: %s ", studentTask.BoundaryHandling)
	log.Printf("Formula: %s ", evaluation.Formula)

}
