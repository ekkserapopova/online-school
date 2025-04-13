package http

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"onlineschool/internal/services/courses"
	"onlineschool/pkg/responser"
	"strconv"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Usecase courses.Usecase
}

type Handler struct {
	logger  *slog.Logger
	usecase courses.Usecase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:  params.Logger,
		usecase: params.Usecase,
	}
}

func (h *Handler) GetCourses(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	courses, err := h.usecase.GetCourses(ctx)
	if err != nil {
		h.logger.Error("Failed to get courses", "error", err.Error())
		responser.SendErr(w, http.StatusInternalServerError, "failed to get courses")
		return
	}

	responser.SendOk(w, http.StatusOK, courses)
}

func (h *Handler) GetCourseById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	h.logger.Info("GetCourseById", "id", vars["id"])

	courseId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Failed to get course id", "error", err.Error())
	}

	course, err := h.usecase.GetCourse(ctx, courseId)
	if err != nil {
		h.logger.Error("Failed to get course", "error", err.Error())
	}

	responser.SendOk(w, http.StatusOK, course)
}

func (h *Handler) GetUsersCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value("userID")
	userID, ok := userIDValue.(int)
	if !ok {
		responser.SendErr(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	courses, err := h.usecase.GetUsersCourses(ctx, userID)
	if err != nil {
		responser.SendErr(w, http.StatusInternalServerError, "Failed to get user courses")
		return
	}

	responser.SendOk(w, http.StatusOK, courses)
}
