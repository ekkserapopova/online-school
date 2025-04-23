package main

import (
	"context"
	"log"
	"onlineschool/internal/app/dsn"
	"onlineschool/internal/app/handler"
	repo "onlineschool/internal/app/repository/postgres"
	"onlineschool/internal/llm/provider"
	"onlineschool/internal/llm/service"
	"onlineschool/internal/redis"
	"os"

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

	repo, _ := repo.NewRepo(connection)

	redisConfig := redis.InitRedisConfig()

	redisClient, err := redis.NewRedisClient(context.Background(), redisConfig)
	if err != nil {
	}

	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY не установлен в .env")
	}

	openAIProvider := provider.NewOpenAIProvider(apiKey, os.Getenv("OPENAI_MODEL"))

	codeEvaluator := service.NewCodeEvaluator(openAIProvider)

	handler := handler.NewHandler(repo, redisClient, *codeEvaluator)

	r := handler.InitRoutes()
	r.Run()
}

func NewRedisClient() (any, any) {
	panic("unimplemented")
}
