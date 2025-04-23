package handler

import (
	"onlineschool/internal/app"
	"onlineschool/internal/llm/service"
	"onlineschool/internal/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo          app.Repo
	redisClient   redis.Client
	codeEvaluator *service.CodeEvaluator
}

func NewHandler(repo app.Repo, redisClient redis.Client, codeEvaluator service.CodeEvaluator) *Handler {
	return &Handler{repo: repo, redisClient: redisClient, codeEvaluator: &codeEvaluator}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{

		//user
		api.POST("/login", h.Login)
		api.POST("/signup", h.Register)
		api.GET("/user/:id", h.GetUserByID)

		// course
		api.GET("/courses", h.GetCourses)
		api.GET("/courses/user", TokenAuth(), h.GetStudentsCourses)
		api.GET("/course/user/:courseID", TokenAuth(), h.GetStudentsCourse)
		api.GET("/course/:id", h.GetCourse)
		api.POST("/enroll/:courseID", TokenAuth(), h.EnrollStudent)
		api.GET("/enrolled/course/:courseID", TokenAuth(), h.IsStudentEnrolledInCourse)

		// lesson
		api.GET("/schedule", TokenAuth(), h.GetLessons)
		api.GET("/lesson/:id", TokenAuth(), h.GetLesson)
		api.GET("/lessons/course/:courseID", TokenAuth(), h.GetLessonsByCourseID)

		//language
		api.GET("/languages", h.GetLanguages)

		//material
		api.GET("/lesson/:id/materials", h.GetMaterials)

		//payment
		api.GET("/payment/course/:courseID", TokenAuth(), h.GetPayment)
		api.POST("/payment/course/:courseID", TokenAuth(), h.AddPayment)
		api.PUT("/payment/course/:courseID", TokenAuth(), h.UpdatePayment)

		//test
		api.GET("/course/:id/tests", TokenAuth(), h.GetTests)
		api.GET("/questions/test/:testID", TokenAuth(), h.GetQuestionsForTest)
		api.POST("/question/:questionID", TokenAuth(), h.AddAnswerByStudent)
		api.GET("/course/:id/test/:testid", TokenAuth(), h.GetTest)
		api.GET("/test/:testID/result", TokenAuth(), h.GetPointsOfTest)
		api.GET("/test/:testID/rightAnswers", TokenAuth(), h.GetRightAnswers)
		api.PUT("/test/:testID/finish", TokenAuth(), h.FinishTest)
		api.GET("/test/:testID/studentAnswers", TokenAuth(), h.GetStudentAnswers)
		//TODO: поменять ручку
		api.GET("/test/:testID/finish", TokenAuth(), h.GetCompletedTest)

		//tasks
		api.GET("/task/:taskID", TokenAuth(), h.GetTask)
		api.POST("/task/:taskID/answer", TokenAuth(), h.AddStudentsTask)

	}

	return r
}
