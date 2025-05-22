package repo

import (
	"errors"
	"log"
	"onlineschool/internal/models"
	"time"
)

func (r *Repo) GetLessons(studentID int, period string) ([]models.LessonResponse, error) {
	var courses []models.Course

	err := r.db.
		Table("students_courses").
		Select("courses.*").
		Joins("join courses on students_courses.course_id = courses.id").
		Where("students_courses.user_id = ?", studentID).
		Find(&courses).Error

	if err != nil {
		return []models.LessonResponse{}, err
	}

	log.Printf("Найдено курсов в расписании: %d", len(courses))

	var result []models.LessonResponse

	// Для каждого курса в расписании загружаем уроки с фильтрацией по дате
	for _, course := range courses {
		var courseModules []models.Module
		err = r.db.Where("course_id = ?", course.ID).
			Preload("Lessons").
			Find(&courseModules).Error

		if err != nil {
			return []models.LessonResponse{}, err
		}

		var filteredLessons []models.Lesson

		for _, module := range courseModules {
			var lessons []models.Lesson

			switch period {
			case "future": // будущие уроки
				lessons, err = r.GetFutureLessons(module.ID)
			case "past": // прошедшие уроки
				lessons, err = r.GetPastLessons(module.ID)
			case "all":
				lessons, err = r.GetAllLessonsByModule(module.ID)
			default:
				return []models.LessonResponse{}, errors.New("Period is not valid")
			}

			if err != nil {
				return []models.LessonResponse{}, err
			}

			filteredLessons = append(filteredLessons, lessons...)
		}

		// Создаём LessonResponse для каждого урока
		for _, lesson := range filteredLessons {
			result = append(result, models.LessonResponse{
				Lesson:     lesson,
				CourseName: course.Name,
			})
		}
	}

	return result, nil
}

func (r *Repo) GetAllLessonsByModule(moduleID int) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("module_id = ?", moduleID).Find(&lessons).Error

	if err != nil {
		return []models.Lesson{}, err
	}
	return lessons, nil
}

func (r *Repo) GetPastLessons(moduleID int) ([]models.Lesson, error) {
	var pastLessons []models.Lesson

	err := r.db.Where("module_id = ? AND start < ?", moduleID, time.Now()).
		Find(&pastLessons).Error

	if err != nil {
		return []models.Lesson{}, err
	}

	return pastLessons, nil
}

func (r *Repo) GetFutureLessons(moduleID int) ([]models.Lesson, error) {
	var futureLessons []models.Lesson

	err := r.db.Where("module_id = ? AND start > ?", moduleID, time.Now()).
		Find(&futureLessons).Error

	if err != nil {
		return []models.Lesson{}, err
	}

	return futureLessons, nil
}

func (r *Repo) GetLesson(lessonId int) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Preload("Materials").Where("id = ?", lessonId).First(&lesson).Error
	if err != nil {
		return lesson, err
	}

	return lesson, nil
}

func (r *Repo) GetLessonsWithoutCourses() ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("module_id = ?", nil).Find(&lessons).Error

	return []models.Lesson{}, err
}

func (r *Repo) AddLesson(lesson models.Lesson) (models.Lesson, error) {
	err := r.db.Create(&lesson).Error
	return lesson, err
}

func (r *Repo) AddLessonToModule(lesson *models.Lesson, moduleID int) error {
	lesson.ModuleID = moduleID
	err := r.db.Updates(&lesson).Where("id = ?", lesson.ID).Error
	return err
}

func (r *Repo) DeleteLesson(lessonId int) error {
	err := r.db.Delete(&models.Lesson{}, "id = ?", lessonId).Error
	return err
}

func (r *Repo) UpdateLesson(lesson models.Lesson) (models.Lesson, error) {
	//var lesson models.Lesson
	err := r.db.Save(&lesson).Error
	return lesson, err
}
