package main

import (
	"log"
	"onlineschool/internal/app/dsn"
	"onlineschool/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
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
	)
	if err != nil {
		panic("cant migrate db")
	}
}
