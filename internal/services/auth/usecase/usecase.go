package usecase

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"log/slog"
	"onlineschool/internal/models"
	"onlineschool/internal/pkg/hasher"
	"onlineschool/internal/services/auth"
	//"onlineschool/internal/models"
	"onlineschool/internal/pkg/jwter"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   auth.Repository
	JWTer  *jwter.JWTer
}

type Usecase struct {
	log  *slog.Logger
	repo auth.Repository
	jwt  *jwter.JWTer
}

func NewUsecase(p Params) *Usecase { // Возвращаем интерфейс auth.Usecase
	return &Usecase{
		log:  p.Logger,
		repo: p.Repo,
		jwt:  p.JWTer,
	}
}

func (uc *Usecase) Login(ctx context.Context, userData *models.User) (string, error) {
	uc.log = uc.log.With("reg", "auth.Usecase.Login")
	user, err := uc.repo.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return "", err
	}

	uc.log.Info("user", "user", user.Email, user.Password, user.Salt)
	uc.log.Info("input pass", "password", userData.Password)

	if !hasher.CheckPassword(userData.Password, user.Salt, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := uc.jwt.GenerateJWT(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *Usecase) Register(ctx context.Context, userData *models.User) (string, error) {

	uc.log = uc.log.With("reg", "auth.Usecase.Register")

	newUser, err := uc.repo.CreateUser(ctx, userData)
	if err != nil {
		return "", err
	}

	token, err := uc.jwt.GenerateJWT(newUser.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
