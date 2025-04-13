package http

import (
	"errors"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"onlineschool/internal/models"
	"onlineschool/internal/services/auth"
	"onlineschool/pkg/reader"
	"onlineschool/pkg/responser"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Usecase auth.Usecase
}

type Handler struct {
	logger  *slog.Logger
	usecase auth.Usecase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:  params.Logger,
		usecase: params.Usecase,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	userData := &models.User{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
		return
	}

	token, err := h.usecase.Login(r.Context(), userData)
	h.logger.Info("input pass HTTP", "data", userData)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrIncorrectData):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectData.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	responser.SendOk(w, http.StatusOK, token)
}

//func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
//	userData := &models.User{}
//	if err := reader.ReadRequestData(r, userData); err != nil {
//		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
//		return
//	}
//
//	token, err := h.usecase.Register(r.Context(), userData)
//	if err != nil {
//		switch {
//		case errors.Is(err, auth.ErrAlreadyExists):
//			responser.SendErr(w, http.StatusBadRequest, auth.ErrAlreadyExists.Error())
//			return
//		case errors.Is(err, auth.ErrIncorrectRole):
//			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectRole.Error())
//			return
//		default:
//			responser.SendErr(w, http.StatusInternalServerError, "internal server error")
//			return
//		}
//	}
//	responser.SendOk(w, http.StatusCreated, token)
//}
