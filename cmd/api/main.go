package main

import (
	"context"
	"log/slog"
	"net/http"
	"onlineschool/internal/app/dsn"
	"onlineschool/internal/pkg/jwter"
	authMiddleware "onlineschool/internal/pkg/middleware" // Добавлен импорт middleware
	"onlineschool/internal/pkg/server"
	"onlineschool/internal/redis"
	"onlineschool/internal/services/auth"
	authHttp "onlineschool/internal/services/auth/delivery/http"
	authRepo "onlineschool/internal/services/auth/repo"
	authUsecase "onlineschool/internal/services/auth/usecase"
	// Импортируем модули курсов
	"onlineschool/internal/services/courses"
	coursesHttp "onlineschool/internal/services/courses/delivery/http"
	coursesRepo "onlineschool/internal/services/courses/repo"
	coursesUsecase "onlineschool/internal/services/courses/usecase"
	"os"
	"os/signal"
	"syscall"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки .env")
		os.Exit(1)
	}

	app := fx.New(
		// Настройка логгера для fx
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		// Предоставляем зависимости
		fx.Provide(
			// Логгер
			func() *slog.Logger {
				return slog.Default()
			},

			// PostgreSQL пул соединений
			func(lc fx.Lifecycle) (*pgxpool.Pool, error) {
				connection := dsn.FromEnv()
				pool, err := pgxpool.New(context.Background(), connection)
				if err != nil {
					return nil, err
				}

				// Управление жизненным циклом пула соединений
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						pool.Close()
						return nil
					},
				})

				return pool, nil
			},

			// SQL Builder
			func() squirrel.StatementBuilderType {
				return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
			},

			// Redis клиент
			func(lc fx.Lifecycle) (*redis.RedisClient, error) {
				redisConfig := redis.InitRedisConfig()
				redisClient, err := redis.NewRedisClient(context.Background(), redisConfig)
				if err != nil {
					return nil, err
				}

				// Управление жизненным циклом Redis клиента
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return redisClient.Close()
					},
				})

				return redisClient, nil
			},

			// JWTer
			jwter.New,

			// --- Auth Service ---
			// Репозиторий Auth
			authRepo.NewRepository,
			// Преобразование *repo.Repository в auth.Repository
			func(r *authRepo.Repository) auth.Repository {
				return r
			},
			// Юзкейс Auth
			authUsecase.NewUsecase,
			// Преобразование *usecase.Usecase в auth.Usecase
			func(u *authUsecase.Usecase) auth.Usecase {
				return u
			},
			// Хендлер Auth
			authHttp.NewHandler,
			// Auth Middleware
			authMiddleware.NewAuthMiddleware, // Добавлен провайдер для middleware

			// --- Courses Service ---
			// Репозиторий Courses
			coursesRepo.NewRepository,
			// Преобразование *repo.Repository в courses.Repository
			func(r *coursesRepo.Repository) courses.Repository {
				return r
			},
			// Юзкейс Courses
			coursesUsecase.NewUsecase,
			// Преобразование *usecase.Usecase в courses.Usecase
			func(u *coursesUsecase.Usecase) courses.Usecase {
				return u
			},
			// Хендлер Courses
			coursesHttp.NewHandler,

			// Роутер
			server.NewRouter,
		),

		// Вызов функций при запуске
		fx.Invoke(
			func(router *server.Router, lc fx.Lifecycle, logger *slog.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Info("Запуск сервера...")
						go func() {
							if err := http.ListenAndServe(":8080", router.Handler); err != nil {
								logger.Error("Ошибка запуска сервера", "error", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Остановка сервера...")
						return nil
					},
				})
			},
		),
	)

	// Управление остановкой приложения
	ctx := context.Background()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if err := app.Start(ctx); err != nil {
		slog.Error("Ошибка запуска приложения", "error", err)
		os.Exit(1)
	}

	<-stop
	app.Stop(ctx)
}
