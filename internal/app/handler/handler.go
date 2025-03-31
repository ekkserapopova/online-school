package handler

import (
	"onlineschool/internal/app"
	"onlineschool/internal/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo        app.Repo
	redisClient redis.Client
}

func NewHandler(repo app.Repo, redisClient redis.Client) *Handler {
	return &Handler{repo: repo, redisClient: redisClient}
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
		// api.GET("/teachers", h.GetTeachers)
		api.GET("/schedule", TokenAuth(), h.GetLessons)
		api.GET("/lesson/:id", h.GetLesson)
		api.GET("/courses", h.GetCourses)
		api.GET("/languages", h.GetLanguages)
		api.GET("/course/:id", h.GetCourse)
		// api.GET("/student/:id", h.GetStudent)
		api.GET("/lesson/:id/materials", h.GetMaterials)
		api.GET("/course/:id/tests", h.GetTests)
		api.GET("/course/:id/test/:testid", h.GetTest)
		api.POST("/login", h.Login)
		api.POST("/signup", h.Register)
		api.GET("/user/:id", h.GetUserByID)
		api.POST("/enroll/:courseID", TokenAuth(), h.EnrollStudent)
		api.GET("/lessons/course/:courseID", h.GetLessonsByCourseID)
	}

	return r
}
