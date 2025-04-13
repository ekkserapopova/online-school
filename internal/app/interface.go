package app

import (
	"onlineschool/internal/models"
)

type Repo interface {
	// GetTeachersList(name string) ([]models.Teacher, error)
	GetCourses(name string) ([]models.Course, error)
	GetUsersCourses(userID int) ([]models.Course, error)
	GetCourse(id int) (models.Course, error)
	GetLessons(studentID int, period string) ([]models.Lesson, error)
	GetLesson(lessonID int) (models.Lesson, error)
	GetTests(courseID int) ([]models.Test, error)
	GetTest(testID int) (models.Test, error)
	GetMaterials(lessonId int) ([]models.Material, error)
	// GetStudent(studentID int) (models.Student, error)

	CreateUser(user *models.User) error

	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) error
	GetByID(id int) (*models.User, error)
	// FindStudentByEmail(email string) (*models.Student, error)

	GetLanguages(name string) ([]models.Language, error)
	// GetPayment()

	EnrollStudent(userID, courseID int) error
	GetLessonsByCourseID(courseID int) ([]models.Lesson, error)

	IsStudentEnrolledInCourse(studentID, courseID int) (bool, error)

	AddPayment(userID, courseID int) (models.Payment, error)
	GetPayment(userID, courseID int) (models.Payment, error)
	UpdatePaymentsStatus(payment *models.Payment) error
}
