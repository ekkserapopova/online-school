package app

import (
	"onlineschool/internal/models"
)

type Repo interface {
	//users
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) error
	GetByID(id int) (*models.User, error)

	//courses
	GetCourses(name string) ([]models.Course, error)
	GetStudentsCourses(userID int) ([]models.Course, error)
	GetStudentsCourse(userID, courseID int) (models.Course, error)
	GetCourse(id int) (models.Course, error)
	EnrollStudent(userID, courseID int) error
	IsStudentEnrolledInCourse(studentID, courseID int) (bool, error)

	//lessons
	GetLessons(studentID int, period string) ([]models.LessonResponse, error)
	GetLesson(lessonID int) (models.Lesson, error)
	GetLessonsByCourseID(courseID int) ([]models.Lesson, error)

	//materials
	GetMaterials(lessonId int) ([]models.Material, error)

	//languages
	GetLanguages(name string) ([]models.Language, error)

	//payments
	AddPayment(userID, courseID int) (models.Payment, error)
	GetPayment(userID, courseID int) (models.Payment, error)
	UpdatePaymentsStatus(payment *models.Payment) error

	//tests
	GetTests(courseID int) ([]models.Test, error)
	GetTest(testID int) (models.Test, error)
	GetQuestionsForTest(testID int) ([]models.Question, error)
	AddAnswerByStudent(studentAnswer models.StudentAnswer) (models.StudentAnswer, error)
	GetPointsOfTest(testID, student_id int) (int, error)
	GetRightAnswers(testID int) ([]models.AnswerResponse, error)
	FinishTest(testID int, studentID int) error
	GetStudentAnswers(testID int, studentID int) ([]models.StudentAnswer, error)
	GetCompletedTest(testID int, studentID int) (models.CompletedTest, error)

	//tasks
	GetTask(taskID int) (*models.Task, error)
}
