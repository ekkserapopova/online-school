package auth

import (
	"context"
	"onlineschool/internal/models"
)

type Usecase interface {
	Login(ctx context.Context, userData *models.User) (string, error)
	Register(ctx context.Context, userData *models.User) (string, error)
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}
