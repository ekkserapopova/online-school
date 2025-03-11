package app

import (
	"onlineschool/internal/models"
	"time"
)

type Repo interface {
	GetTeachersList(name string) ([]models.Teacher, error)
	GetCourses(name string) ([]models.Course, error)
	GetCourse(id int) (models.Course, error)
	GetSchedule(studentID int, from, to time.Time) ([]models.Lesson, error)
	GetLesson(lessonID int) (models.Lesson, error)
	GetTests(courseID int) ([]models.Test, error)
	GetTest(testID int) (models.Test, error)
	GetMaterials(lessonId int) ([]models.Material, error)
	GetStudent(studentID int) (models.Student, error)
	// GetPayment()
}
