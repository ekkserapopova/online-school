package main

import (
	"log"
	"onlineschool/internal/app/dsn"
	"onlineschool/internal/app/handler"
	repo "onlineschool/internal/app/repository"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	connection := dsn.FromEnv()

	repo, _ := repo.NewRepository(connection)

	handler := handler.NewHandler(repo)

	r := handler.InitRoutes()
	r.Run()
}
