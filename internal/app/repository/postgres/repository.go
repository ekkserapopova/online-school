package repo

import (
	"log"
	"onlineschool/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(connectionString string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Language{},
		&models.Course{},
		&models.Lesson{},
		&models.Module{},
		&models.Payment{},
		&models.Homework{},
		&models.Material{},
		&models.Review{},
		&models.Test{},
		&models.Question{},
		&models.AnswerVariant{},
		&models.StudentAnswer{},
		&models.CompletedTest{},
		&models.Task{},
		&models.StudentTask{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return &Repo{
		db: db,
	}, nil
}
