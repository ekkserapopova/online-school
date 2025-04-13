package repo

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log/slog"
	"onlineschool/internal/models"
	"onlineschool/internal/pkg/hasher"
	"onlineschool/internal/services/auth"
)

type Params struct {
	fx.In

	Pool    *pgxpool.Pool
	Logger  *slog.Logger
	Builder squirrel.StatementBuilderType
}

type Repository struct {
	pool    *pgxpool.Pool
	log     *slog.Logger
	builder squirrel.StatementBuilderType
}

func NewRepository(params Params) *Repository {
	return &Repository{
		pool:    params.Pool,
		log:     params.Logger,
		builder: params.Builder,
	}
}

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := repo.builder.
		Select("id", "email", "password", "salt").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	user := &models.User{}

	err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Salt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Debug("get user by email: user not found")
			return nil, errors.New("user not found")
		}
		repo.log.Error("failed to get user: " + err.Error())
		return nil, err
	}

	return user, nil
}

func (repo *Repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	if user.Name == "" {
		return nil, errors.New("name is required")
	}

	if user.Surname == "" {
		return nil, errors.New("surname is required")
	}

	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	if user.Phone == "" {
		return nil, errors.New("phone is required")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

	salt, err := hasher.GenerateSalt()
	if err != nil {
		repo.log.Error("failed to generate salt: " + err.Error())
		return nil, err
	}

	hashPassword := hasher.HashPassword(user.Password, salt)

	// Создаем билдер для вставки
	insertBuilder := repo.builder.
		Insert("users").
		Suffix("RETURNING id, email")

	// Создаем слайсы для динамического добавления колонок и значений
	columns := []string{}
	values := []interface{}{}

	// Добавляем только непустые поля
	// Обязательные поля
	columns = append(columns, "name", "surname", "email", "phone", "password", "salt", "is_admin", "is_active")
	values = append(values, user.Name, user.Surname, user.Email, user.Phone, hashPassword, salt, false, true)

	if !user.Birth.IsZero() {
		columns = append(columns, "birth")
		values = append(values, user.Birth)
	}

	if user.Photo != "" {
		columns = append(columns, "photo")
		values = append(values, user.Photo)
	}

	// Добавляем колонки и значения к билдеру
	insertBuilder = insertBuilder.Columns(columns...).Values(values...)

	// Строим SQL-запрос
	query, args, err := insertBuilder.ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	createdUser := &models.User{}
	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Password,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			repo.log.Warn("user already exists")
			return nil, auth.ErrAlreadyExists
		}
		repo.log.Error("failed to create user: " + err.Error())
		return nil, err
	}

	return createdUser, nil
}
