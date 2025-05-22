package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineschool/internal/models"
	"strconv"
	"time"
)

func (h *Handler) GetTest(c *gin.Context) {
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

func (h *Handler) RandomGenerateTest(c *gin.Context) {
	//TODO: проверка на роль
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

	test, err := h.repo.RandomGenerateTest(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"test": test})
}

func (h *Handler) CreateTest(c *gin.Context) {
	//TODO: проверка на роль
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	input := models.Test{}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ModuleID = moduleID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	test, err := h.repo.CreateTest(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"test": test})
}

func (h *Handler) DeleteTest(c *gin.Context) {
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

	err = h.repo.DeleteTest(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

//func (h *Handler) UpdateTest(c *gin.Context) {
//	_, ok := c.Get("userId")
//	if !ok {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//		return
//	}
//
//	testID, err := strconv.Atoi(c.Param("testID"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test id"})
//		return
//
//	}
//
//	test, err := h.repo.GetTest(testID)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	input := models.Test{}
//	if err := c.ShouldBindJSON(&input); err != nil {
//		log.Println("Error binding json:", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	test.UpdatedAt = time.Now()
//	test.Name = input.Name
//	test.Deadline = input.Deadline
//	test.Description = input.Description
//
//	err = h.repo.UpdateTest(test)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"test": test})
//}

func (h *Handler) GetQuestionsForModule(c *gin.Context) {
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

	if teacher.IsTeacher == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	questions, err := h.repo.GetQuestionsForModule(moduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *Handler) AddQuestion(c *gin.Context) {
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

	if teacher.IsTeacher == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	moduleID, err := strconv.Atoi(c.Param("moduleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	input := models.Question{}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	input.ModuleID = moduleID

	question, err := h.repo.AddQuestion(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, gin.H{"question": question})
}

func (h *Handler) AddAnswer(c *gin.Context) {
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

	if teacher.IsTeacher == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	input := models.AnswerVariant{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	input.QuestionID = questionID

	answerVariant, err := h.repo.AddAnswer(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"answer_variant": answerVariant})

}

func (h *Handler) UpdateTest(c *gin.Context) {
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

	if teacher.IsTeacher == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	test, err := h.repo.GetTest(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if name, ok := input["name"].(string); ok {
		test.Name = name
	}
	if desc, ok := input["description"].(string); ok {
		test.Description = desc
	}

	if deadlineStr, ok := input["deadline"].(string); ok {
		deadline, err := time.Parse(time.RFC3339, deadlineStr)
		if err == nil {
			test.Deadline = deadline
		}
	}

	if countQuestions, ok := input["count_questions"].(float64); ok {
		test.CountQuestions = int(countQuestions)
	}

	if timeLimit, ok := input["time_limit"].(float64); ok {
		test.TimeLimit = int(timeLimit)
	}

	if is_active, ok := input["is_active"].(bool); ok {
		test.IsActive = is_active
	}
	test.UpdatedAt = time.Now()

	err = h.repo.UpdateTest(test)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"test": test})
}

func (h *Handler) CreateStudentTest(c *gin.Context) {
	studentID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	test, err := h.repo.CreateStudentTest(testID, studentID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"test": test})
}

func (h *Handler) GetStudentsTests(c *gin.Context) {
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

	if teacher.IsTeacher == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	testID, err := strconv.Atoi(c.Param("testID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}
	completedTests, err := h.repo.GetStudentsTests(testID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students_tests": completedTests})
}

func (h *Handler) DeleteAnswer(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	answerID, err := strconv.Atoi(c.Param("answerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	err = h.repo.DeleteAnswer(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module id"})
		return
	}

	err = h.repo.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question id"})
		return
	}

	question, err := h.repo.GetQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := make(map[string]interface{})
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if text, ok := input["text"].(string); ok {
		question.Text = text
	}

	err = h.repo.UpdateQuestion(question)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"question": question})
}

func (h *Handler) UpdateAnswerVariant(c *gin.Context) {
	_, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	answerID, err := strconv.Atoi(c.Param("answerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer id"})
		return
	}

	answer, err := h.repo.GetAnswerVariant(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := make(map[string]interface{})
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if text, ok := input["text"].(string); ok {
		answer.Text = text
	}

	err = h.repo.UpdateAnswerVariant(answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
