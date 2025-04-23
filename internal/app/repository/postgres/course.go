package repo

import (
	"gorm.io/gorm"
	"onlineschool/internal/models"
)

func (r *Repo) GetCourses(name string) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&courses).Error
	if err != nil {
		return courses, err
	}
	return courses, nil
}

func (r *Repo) GetCourse(courseID int) (models.Course, error) {
	var course models.Course

	err := r.db.Preload("Modules", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Where("id = ?", courseID).
		First(&course).Error

	if err != nil {
		return course, err
	}

	for i, module := range course.Modules {
		if module.OpenForAll {
			var lessons []models.Lesson
			if err := r.db.Where("module_id = ?", module.ID).
				Order("id ASC").
				Find(&lessons).Error; err != nil {
				return course, err
			}
			course.Modules[i].Lessons = lessons
		}
	}

	return course, nil
}

func (r *Repo) GetStudentsCourses(userID int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Table("students_courses").
		Select("courses.*").
		Joins("join courses on students_courses.course_id = courses.id").
		Where("students_courses.user_id = ?", userID).
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *Repo) GetLanguages(name string) ([]models.Language, error) {
	var languages []models.Language
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&languages).Error
	if err != nil {
		return languages, err
	}
	return languages, nil
}

func (r *Repo) GetStudentsCourse(userID, courseID int) (models.Course, error) {
	var course models.Course
	err := r.db.Table("students_courses").
		Select("courses.*").Preload("Modules.Lessons").Preload("Modules.Tests").Preload("Modules.Tasks").
		Joins("join courses on students_courses.course_id = courses.id").
		Where("students_courses.course_id = ? AND students_courses.user_id = ?", courseID, userID).
		Take(&course).Error

	if err != nil {
		return course, err
	}

	return course, nil
}

func (r *Repo) EnrollStudent(userID, courseID int) error {

	err := r.db.Table("students_courses").Create(map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetLessonsByCourseID(courseID int) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("course_id = ?", courseID).Find(&lessons).Error
	if err != nil {
		return lessons, err
	}

	return lessons, nil
}

func (r *Repo) IsStudentEnrolledInCourse(studentID, courseID int) (bool, error) {
	var count int64
	err := r.db.Table("students_courses").
		Where("user_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
