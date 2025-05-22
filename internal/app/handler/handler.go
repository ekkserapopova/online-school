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
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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
		api.GET("/user/photo", TokenAuth(), h.GetPhoto)
		api.POST("/user/photo", TokenAuth(), h.AddPhoto)

		// course

		courses := api.Group("/courses")
		{
			courses.POST("", TokenAuth(), h.CreateCourse)
			courses.GET("", h.GetCourses)
			courses.GET("/user", TokenAuth(), h.GetStudentsCourses)
			courses.GET("/progress", TokenAuth(), h.GetUserProgressForAllCourses)
			courses.GET("/teacher", TokenAuth(), h.GetTeachersCourses)

			course := courses.Group("/:courseID")
			{
				course.GET("", h.GetCourse)
				course.GET("/image", h.GetCoursePhoto)
				course.PATCH("/image", TokenAuth(), h.UpdateCoursePhoto)
				course.GET("/user", TokenAuth(), h.GetStudentsCourse)
				course.GET("/teacher", TokenAuth(), h.GetTeachersCourse)
				course.POST("/enroll", TokenAuth(), h.EnrollStudent)
				course.GET("/enrolled", TokenAuth(), h.IsStudentEnrolledInCourse)

			}

		}

		//module
		api.POST("/courses/:courseID/modules", TokenAuth(), h.AddModuleToCourse)
		// lesson
		lessons := api.Group("/lessons")
		{
			lessons.GET("", TokenAuth(), h.GetLessons)

			lessons.GET("/:lessonID", TokenAuth(), h.GetLesson)
			lessons.DELETE("/:lessonID", TokenAuth(), h.DeleteLesson)
			lessons.PATCH("/:lessonID", TokenAuth(), h.UpdateLesson)

			lessons.GET("/course/:courseID", TokenAuth(), h.GetLessonsByCourseID)
		}

		api.GET("/modules/:moduleID", TokenAuth(), h.GetModule)
		//api.POST("/module/:moduleID/lesson", TokenAuth(), h.AddLesson)

		//language
		api.GET("/languages", h.GetLanguages)

		//material
		api.GET("/lessons/:lessonID/materials", h.GetMaterials)

		//payment
		api.GET("/payment/course/:courseID", TokenAuth(), h.GetPayment)
		api.POST("/payment/course/:courseID", TokenAuth(), h.AddPayment)
		api.PUT("/payment/course/:courseID", TokenAuth(), h.UpdatePayment)

		//test
		tests := api.Group("/tests")
		{
			tests.GET("/:testID/questions", TokenAuth(), h.GetQuestionsForTest)
			tests.GET("/:testID/result", TokenAuth(), h.GetPointsOfTest)
			tests.GET("/:testID/results", TokenAuth(), h.GetStudentsTests)
			tests.GET("/:testID/rightAnswers", TokenAuth(), h.GetRightAnswers)
			tests.PUT("/:testID/finish", TokenAuth(), h.FinishTest)
			tests.GET("/:testID/studentAnswers", TokenAuth(), h.GetStudentAnswers)
			//TODO: поменять ручку
			tests.GET("/:testID/finish", TokenAuth(), h.GetCompletedTest)
			tests.POST("/:testID/generate-questions", TokenAuth(), h.RandomGenerateTest)
			tests.DELETE("/:testID", TokenAuth(), h.DeleteTest)
			tests.GET("/:testID", TokenAuth(), h.GetTest)
			tests.POST("/:testID/student", TokenAuth(), h.CreateStudentTest)

		}

		api.GET("/materials/:materialID", h.GetMaterialPDF)

		api.GET("/course/:courseID/tests", TokenAuth(), h.GetTests)
		api.POST("/question/:questionID", TokenAuth(), h.AddAnswerByStudent)
		api.GET("/courses/:courseID/test/:testID", TokenAuth(), h.GetTest)

		api.POST("/module/:moduleID/test", TokenAuth(), h.CreateTest)

		//api.POST("/courses/:courseID/module", TokenAuth(), h.AddModuleToCourse)

		//tasks

		tasks := api.Group("/tasks")
		{
			tasks.GET("/:taskID", TokenAuth(), h.GetTask)
			tasks.GET("", TokenAuth(), h.GetAllTasks)
			tasks.GET("/students/:taskID", TokenAuth(), h.GetTaskAnswersForTeacher)
			tasks.POST("/:taskID/answer", TokenAuth(), h.AddStudentsTask)
			tasks.GET("/:taskID/answers", TokenAuth(), h.GetStudentTasks)
			tasks.GET("/:taskID/answer", TokenAuth(), h.GetStudentTask)
			tasks.GET("/:taskID/score", TokenAuth(), h.GetFinalScore)
			tasks.DELETE("/:taskID", TokenAuth(), h.DeleteTask)
			tasks.PATCH("/:taskID", TokenAuth(), h.UpdateTask)

		}

		api.POST("/module/:moduleID/task", TokenAuth(), h.AddTask)
		api.POST("/module/:moduleID/lesson", TokenAuth(), h.AddLesson)
		api.GET("/module/:moduleID/questions", TokenAuth(), h.GetQuestionsForModule)

		api.POST("/module/:moduleID/questions", TokenAuth(), h.AddQuestion)
		api.POST("/questions/:questionID/answers", TokenAuth(), h.AddAnswer)
		api.DELETE("/questions/:questionID/answers/:answerID", TokenAuth(), h.DeleteAnswer)
		api.DELETE("/questions/:questionID", TokenAuth(), h.DeleteQuestion)

		api.PATCH("/tests/:testID", TokenAuth(), h.UpdateTest)

		api.DELETE("/modules/:moduleID", TokenAuth(), h.DeleteModule)
		api.PATCH("/modules/:moduleID", TokenAuth(), h.UpdateModule)

		api.POST("lessons/:lessonID/materials", TokenAuth(), h.CreateMaterial)
		api.DELETE("lessons/:lessonID/materials/:materialID", TokenAuth(), h.DeleteMaterial)
		api.PATCH("lessons/:lessonID/materials/:materialID", TokenAuth(), h.UpdateMaterial)
		api.PATCH("/questions/:questionID", TokenAuth(), h.UpdateQuestion)
		api.PATCH("/variants-answers/:answerID", TokenAuth(), h.UpdateAnswerVariant)
	}

	return r
}
