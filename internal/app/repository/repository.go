package repo

import (
	"fmt"
	"log"
	"onlineschool/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepository(connectionString string) (*Repo, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Teacher{},
		&models.Student{},
		&models.Course{},
		&models.Lesson{},
		&models.Payment{},
		&models.Schedule{},
		&models.Homework{},
		&models.Material{},
		&models.Review{},
		&models.Test{},
		&models.Question{},
		&models.Answer{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return &Repo{
		db: db,
	}, nil
}

func (r *Repo) GetTeachersList(name string) ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&teachers).Error
	if err != nil {
		return teachers, err
	}
	return teachers, nil
}

func (r *Repo) GetCourses(name string) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&courses).Error
	if err != nil {
		return courses, err
	}
	return courses, nil
}

func (r *Repo) GetCourse(id int) (models.Course, error) {
	var course models.Course
	err := r.db.Where("id = ?", id).First(&course).Error
	if err != nil {
		return course, err
	}
	return course, nil
}

func (r *Repo) GetSchedule(studentID int, start, end time.Time) ([]models.Lesson, error) {
	var schedule models.Schedule
	err := r.db.Preload("Courses.Lessons").Where("student_id = ?", studentID).First(&schedule).Error

	if err != nil {
		return nil, nil
	}

	var lessons []models.Lesson

	for _, course := range schedule.Courses {
		r.db.Where("start >= ? AND start <= ?", start, end).Find(&lessons)
		fmt.Println(course.Lessons)
	}

	return lessons, nil
}

func (r *Repo) GetMaterials(lessonID int) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.Where("lesson_id = ?", lessonID).Find(&materials).Error
	if err != nil {
		return materials, err
	}

	return materials, nil
}

func (r *Repo) GetStudent(studentID int) (models.Student, error) {
	var student models.Student
	err := r.db.Where("id = ?", studentID).First(&student).Error
	if err != nil {
		return student, err
	}

	return student, nil
}

func (r *Repo) GetLesson(lessonId int) (models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.Where("id = ?", lessonId).First(&lesson).Error
	if err != nil {
		return lesson, err
	}

	return lesson, nil
}

func (r *Repo) GetTest(testID int) (models.Test, error) {
	var test models.Test
	err := r.db.Preload("Questions").Where("id = ?", testID).First(&test).Error
	if err != nil {
		return test, err
	}

	return test, nil
}

func (r *Repo) GetTests(courseID int) ([]models.Test, error) {
	var tests []models.Test
	err := r.db.Preload("Questions").Where("course_id = ?", courseID).Find(&tests).Error
	if err != nil {
		return tests, err
	}

	return tests, nil
}
