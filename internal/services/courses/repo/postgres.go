package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log/slog"
	"onlineschool/internal/models"
	"time"
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

func (repo *Repository) GetCourses(ctx context.Context) ([]*models.Course, error) {
	// Запрос с JOIN для загрузки курсов, преподавателей и уроков
	//TODO: доделать запрос
	query, args, err := repo.builder.
		Select(
			"c.id as course_id",
			"c.name as course_name",
			"c.description as course_description",
			"c.difficulty as course_difficulty",
			"c.price as course_price",
			"c.is_active as course_is_active",
			"c.created_at as course_created_at",
			"c.updated_at as course_updated_at",
			"u.id as teacher_id",
			"u.surname as teacher_surname",
			"u.name as teacher_name",
			"u.email as teacher_email",
			"u.overview as teacher_overview",
			"l.id as lesson_id",
			"l.name as lesson_name",
			"l.description as lesson_description",
			"l.created_at as lesson_created_at",
			"l.updated_at as lesson_updated_at",
		).
		From("courses c").
		LeftJoin("users u ON c.teacher_id = u.id").
		LeftJoin("lessons l ON c.id = l.course_id").
		OrderBy("c.id, l.id"). // Сортировка важна для группировки результатов
		ToSql()

	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		repo.log.Error("failed to query data: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	courses := make([]*models.Course, 0)
	coursesMap := make(map[int]*models.Course)

	// Обработка результатов запроса
	for rows.Next() {
		var courseID, teacherID int
		var lessonID sql.NullInt64 // Может быть NULL, если у курса нет уроков
		var courseName, courseDescription string
		var courseDifficulty, coursePrice int
		var courseIsActive bool
		var courseCreatedAt, courseUpdatedAt time.Time
		var teacherSurname, teacherName, teacherEmail, teacherOverview string
		var lessonName, lessonDescription sql.NullString  // Может быть NULL
		var lessonCreatedAt, lessonUpdatedAt sql.NullTime // Может быть NULL

		err := rows.Scan(
			&courseID,
			&courseName,
			&courseDescription,
			&courseDifficulty,
			&coursePrice,
			&courseIsActive,
			&courseCreatedAt,
			&courseUpdatedAt,
			&teacherID,
			&teacherSurname,
			&teacherName,
			&teacherEmail,
			&teacherOverview,
			&lessonID,
			&lessonName,
			&lessonDescription,
			&lessonCreatedAt,
			&lessonUpdatedAt,
		)

		if err != nil {
			repo.log.Error("failed to scan data: " + err.Error())
			return nil, err
		}

		// Проверяем, есть ли уже такой курс в нашей карте
		course, exists := coursesMap[courseID]
		if !exists {
			// Создаем новый курс
			course = &models.Course{
				ID:          courseID,
				Name:        courseName,
				Description: courseDescription,
				Difficulty:  courseDifficulty,
				Price:       coursePrice,
				TeacherID:   teacherID,
				IsActive:    courseIsActive,
				CreatedAt:   courseCreatedAt,
				UpdatedAt:   courseUpdatedAt,
				Lessons:     make([]models.Lesson, 0),
			}

			// Добавляем курс в слайс и карту
			courses = append(courses, course)
			coursesMap[courseID] = course
		}

		// Если есть данные об уроке, добавляем его к курсу
		if lessonID.Valid {
			lesson := models.Lesson{
				ID:          int(lessonID.Int64),
				Name:        lessonName.String,
				Description: lessonDescription.String,
			}

			if lessonCreatedAt.Valid {
				lesson.CreatedAt = lessonCreatedAt.Time
			}

			if lessonUpdatedAt.Valid {
				lesson.UpdatedAt = lessonUpdatedAt.Time
			}

			course.Lessons = append(course.Lessons, lesson)
		}
	}

	if err := rows.Err(); err != nil {
		repo.log.Error("error during rows iteration: " + err.Error())
		return nil, err
	}

	return courses, nil
}

func (repo *Repository) GetCourse(ctx context.Context, id int) (*models.Course, error) {
	queryBuilder, args, err := repo.builder.
		Select(
			"c.id", "c.name", "c.description", "c.difficulty", "c.price",
			"c.teacher_id", "c.is_active", "c.created_at", "c.updated_at",
			// Для учителя
			"json_build_object('id', t.id, 'name', t.name, 'surname', t.surname, 'email', t.email) AS teacher",
			// Для уроков
			"COALESCE(json_agg(json_build_object('id', l.id, 'name', l.name, 'description', l.description)) FILTER (WHERE l.id IS NOT NULL), '[]') AS lessons",
		).
		From("courses c").
		LeftJoin("users t ON t.id = c.teacher_id").
		LeftJoin("lessons l ON l.course_id = c.id").
		Where(squirrel.Eq{"c.id": id}).
		GroupBy("c.id", "t.id").
		ToSql()

	if err != nil {
		repo.log.Error("failed to build query: " + err.Error())
		return nil, err
	}

	repo.log.Info("Final SQL: " + queryBuilder)

	course := &models.Course{}
	var courseCreatedAt, courseUpdatedAt time.Time
	var lessonsJSON []byte
	var teacherJSON []byte

	err = repo.pool.QueryRow(ctx, queryBuilder, args...).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Difficulty,
		&course.Price,
		&course.TeacherID,
		&course.IsActive,
		&courseCreatedAt,
		&courseUpdatedAt,
		&teacherJSON,
		&lessonsJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Курс не найден
		}
		repo.log.Error("failed to scan course: " + err.Error())
		return nil, err
	}

	course.CreatedAt = courseCreatedAt
	course.UpdatedAt = courseUpdatedAt

	// Десериализуем JSON с данными учителя
	if teacherJSON != nil && len(teacherJSON) > 0 && string(teacherJSON) != "null" {
		if course.Teacher == nil {
			course.Teacher = &models.User{} // Инициализируем, если nil
		}
		err = json.Unmarshal(teacherJSON, &course.Teacher)
		if err != nil {
			repo.log.Error("failed to unmarshal teacher JSON: " + err.Error())
		}
	}

	// Десериализуем JSON с уроками
	if lessonsJSON != nil && len(lessonsJSON) > 0 && string(lessonsJSON) != "null" {
		err = json.Unmarshal(lessonsJSON, &course.Lessons)
		if err != nil {
			repo.log.Error("failed to unmarshal lessons JSON: " + err.Error())
		}
	}

	return course, nil
}

func (repo *Repository) GetUsersCourses(ctx context.Context, userID int) ([]*models.Course, error) {
	// Запрос с JOIN для загрузки курсов, преподавателей и уроков
	//TODO: доделать запрос
	query, args, err := repo.builder.
		Select(
			"c.id as course_id",
			"c.name as course_name",
			"c.description as course_description",
			"c.difficulty as course_difficulty",
			"c.price as course_price",
			"c.is_active as course_is_active",
			"c.created_at as course_created_at",
			"c.updated_at as course_updated_at",
			"u.id as teacher_id",
			"u.surname as teacher_surname",
			"u.name as teacher_name",
			"u.email as teacher_email",
			"u.overview as teacher_overview",
			"l.id as lesson_id",
			"l.name as lesson_name",
			"l.description as lesson_description",
			"l.created_at as lesson_created_at",
			"l.updated_at as lesson_updated_at",
		).
		From("courses c").
		LeftJoin("users u ON c.teacher_id = u.id").
		LeftJoin("lessons l ON c.id = l.course_id").
		Where(squirrel.Eq{"c.user_id": userID}).
		OrderBy("c.id, l.id"). //для группировки результатов
		ToSql()

	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		repo.log.Error("failed to query data: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	courses := make([]*models.Course, 0)
	coursesMap := make(map[int]*models.Course)

	// Обработка результатов запроса
	for rows.Next() {
		var courseID, teacherID int
		var lessonID sql.NullInt64 // Может быть NULL, если у курса нет уроков
		var courseName, courseDescription string
		var courseDifficulty, coursePrice int
		var courseIsActive bool
		var courseCreatedAt, courseUpdatedAt time.Time
		var teacherSurname, teacherName, teacherEmail, teacherOverview string
		var lessonName, lessonDescription sql.NullString  // Может быть NULL
		var lessonCreatedAt, lessonUpdatedAt sql.NullTime // Может быть NULL

		err := rows.Scan(
			&courseID,
			&courseName,
			&courseDescription,
			&courseDifficulty,
			&coursePrice,
			&courseIsActive,
			&courseCreatedAt,
			&courseUpdatedAt,
			&teacherID,
			&teacherSurname,
			&teacherName,
			&teacherEmail,
			&teacherOverview,
			&lessonID,
			&lessonName,
			&lessonDescription,
			&lessonCreatedAt,
			&lessonUpdatedAt,
		)

		if err != nil {
			repo.log.Error("failed to scan data: " + err.Error())
			return nil, err
		}

		// Проверяем, есть ли уже такой курс в нашей карте
		course, exists := coursesMap[courseID]
		if !exists {
			// Создаем новый курс
			course = &models.Course{
				ID:          courseID,
				Name:        courseName,
				Description: courseDescription,
				Difficulty:  courseDifficulty,
				Price:       coursePrice,
				TeacherID:   teacherID,
				IsActive:    courseIsActive,
				CreatedAt:   courseCreatedAt,
				UpdatedAt:   courseUpdatedAt,
				Lessons:     make([]models.Lesson, 0),
			}

			// Добавляем курс в слайс и карту
			courses = append(courses, course)
			coursesMap[courseID] = course
		}

		// Если есть данные об уроке, добавляем его к курсу
		if lessonID.Valid {
			lesson := models.Lesson{
				ID:          int(lessonID.Int64),
				Name:        lessonName.String,
				Description: lessonDescription.String,
			}

			if lessonCreatedAt.Valid {
				lesson.CreatedAt = lessonCreatedAt.Time
			}

			if lessonUpdatedAt.Valid {
				lesson.UpdatedAt = lessonUpdatedAt.Time
			}

			course.Lessons = append(course.Lessons, lesson)
		}
	}

	if err := rows.Err(); err != nil {
		repo.log.Error("error during rows iteration: " + err.Error())
		return nil, err
	}

	return courses, nil
}
