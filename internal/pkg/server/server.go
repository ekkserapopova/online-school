package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	authMiddleware "onlineschool/internal/pkg/middleware" // Импортируем middleware
	authHandler "onlineschool/internal/services/auth/delivery/http"
	courseHandler "onlineschool/internal/services/courses/delivery/http"
)

type RouterParams struct {
	fx.In

	Logger        *slog.Logger
	AuthHandler   *authHandler.Handler
	CourseHandler *courseHandler.Handler
	AuthMD        *authMiddleware.AuthMiddleware // Добавляем middleware в параметры
}

type Router struct {
	Handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	//api.Use(middleware.CORSMiddleware)

	v1 := api.PathPrefix("/v1").Subrouter()
	//v1.HandleFunc("/dummyLogin", p.AuthHandler.DummyLogin).Methods(http.MethodPost, http.MethodOptions)
	//v1.HandleFunc("/register", p.AuthHandler.).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/login", p.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	v1.HandleFunc("/courses", p.CourseHandler.GetCourses).Methods(http.MethodGet, http.MethodOptions)
	v1.HandleFunc("/course/{id}", p.CourseHandler.GetCourseById).Methods(http.MethodGet, http.MethodOptions)

	// Защищаем маршрут с помощью middleware
	v1.Handle("/courses/user", p.AuthMD.AuthMiddleware(http.HandlerFunc(p.CourseHandler.GetUsersCourses))).Methods(http.MethodGet, http.MethodOptions)

	router := &Router{
		Handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
