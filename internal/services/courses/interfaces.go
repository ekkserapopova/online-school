package courses

import (
	"context"
	"onlineschool/internal/models"
)

type Repository interface {
	GetCourses(ctx context.Context) ([]*models.Course, error)
	GetCourse(ctx context.Context, id int) (*models.Course, error)
	GetUsersCourses(ctx context.Context, userID int) ([]*models.Course, error)
}

type Usecase interface {
	GetCourses(ctx context.Context) ([]*models.Course, error)
	GetCourse(ctx context.Context, id int) (*models.Course, error)
	GetUsersCourses(ctx context.Context, userID int) ([]*models.Course, error)
}
