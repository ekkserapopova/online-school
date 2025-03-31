package main

import (
	"context"
	"log"
	"onlineschool/internal/app/dsn"
	"onlineschool/internal/app/handler"
	repo "onlineschool/internal/app/repository"
	"onlineschool/internal/redis"

	"github.com/joho/godotenv"
)

// @title BITOP
// @version 1.0
// @description Bmstu Open IT Platform

// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes https http
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	connection := dsn.FromEnv()

	repo, _ := repo.NewRepository(connection)

	redisConfig := redis.InitRedisConfig()

	redisClient, err := redis.NewRedisClient(context.Background(), redisConfig)
	if err != nil {
	}

	handler := handler.NewHandler(repo, redisClient)

	r := handler.InitRoutes()
	r.Run()
}

func NewRedisClient() (any, any) {
	panic("unimplemented")
}
