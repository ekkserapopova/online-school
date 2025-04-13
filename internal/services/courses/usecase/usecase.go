package usecase

import (
	"context"
	"go.uber.org/fx"
	"log/slog"
	"onlineschool/internal/models"
	"onlineschool/internal/services/courses"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   courses.Repository
}

type Usecase struct {
	log  *slog.Logger
	repo courses.Repository
}

func NewUsecase(p Params) *Usecase {
	return &Usecase{
		log:  p.Logger,
		repo: p.Repo,
	}
}

func (u *Usecase) GetCourses(ctx context.Context) ([]*models.Course, error) {
	u.log.Info("GetCourses")
	courses, err := u.repo.GetCourses(ctx)

	if err != nil {
		u.log.Error("Failed to get courses", "error", err.Error())
		return nil, err
	}

	u.log.Info("Successfully retrieved courses", "count", len(courses))
	return courses, nil
}

func (u *Usecase) GetCourse(ctx context.Context, id int) (*models.Course, error) {
	u.log.Info("GetCourse", "id", id)
	course, err := u.repo.GetCourse(ctx, id)

	if err != nil {
		u.log.Error("Failed to get course", "error", err.Error())
		return nil, err
	}

	u.log.Info("Successfully retrieved course")
	return course, nil

}

func (u *Usecase) GetUsersCourses(ctx context.Context, userID int) ([]*models.Course, error) {
	u.log.Info("GetCourses for user", "id", userID)
	courses, err := u.repo.GetUsersCourses(ctx, userID)

	if err != nil {
		u.log.Error("Failed to get courses", "error", err.Error())
		return nil, err
	}

	u.log.Info("Successfully retrieved courses", "count", len(courses))
	return courses, nil
}
