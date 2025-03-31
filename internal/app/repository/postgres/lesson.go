package repo

import (
	"log"
	"onlineschool/internal/models"
	"sort"
	"time"
)

func (r *Repo) GetLessons(studentID int, period string) ([]models.Lesson, error) {
	var schedule models.Schedule

	// Загружаем расписание с предзагрузкой курсов
	err := r.db.Preload("Courses").Where("student_id = ?", studentID).First(&schedule).Error
	if err != nil {
		log.Printf("Ошибка при загрузке расписания: %v", err)
		return nil, err
	}

	log.Println(schedule.Courses)

	log.Printf("Найдено курсов в расписании: %d", len(schedule.Courses))

	var allLessons []models.Lesson

	// Для каждого курса в расписании загружаем уроки с фильтрацией по дате
	for _, course := range schedule.Courses {
		var courseLessons []models.Lesson
		switch period {
		case "future": // будущие уроки
			err := r.db.Preload("Course.Teacher").Where("course_id = ? and start >= ?",
				course.ID, time.Now()).Find(&courseLessons).Error

			if err != nil {
				log.Printf("Ошибка при загрузке уроков для курса %d: %v", course.ID, err)
				continue
			}

			log.Printf("Курс %d (%s): найдено уроков: %d", course.ID, course.Name, len(courseLessons))
		case "past": // прошедшие уроки
			err := r.db.Preload("Course.Teacher").Where("course_id = ? and start < ?",
				course.ID, time.Now()).Find(&courseLessons).Error

			if err != nil {
				log.Printf("Ошибка при загрузке уроков для курса %d: %v", course.ID, err)
				continue
			}
		case "all":
			err := r.db.Preload("Course.Teacher").Where("course_id = ?",
				course.ID).Find(&courseLessons).Error
			if err != nil {
				log.Printf("Ошибка при загрузке уроков для курса %d: %v", course.ID, err)
				continue
			}

		}
		allLessons = append(allLessons, courseLessons...)
	}
	sort.Slice(allLessons, func(i, j int) bool {
		return allLessons[i].Start.Before(allLessons[j].Start)
	})
	log.Printf("Всего найдено уроков: %d", len(allLessons))
	return allLessons, nil
}

func (r *Repo) GetLesson(lessonId int) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Where("id = ?", lessonId).First(&lesson).Error
	if err != nil {
		return lesson, err
	}

	return lesson, nil
}
