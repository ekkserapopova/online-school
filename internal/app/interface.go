package app

import (
	repo "onlineschool/internal/app/repository/postgres"
	"onlineschool/internal/models"
)

//type UserRepository interface {
//	AddPhoto(user *models.User, path string) error
//	CreateUser(user *models.User) error
//	FindByEmail(email string) (*models.User, error)
//	FindByPhone(phone string) error
//	GetByID(id int) (*models.User, error)
//}

type Repo interface {
	//users
	AddPhoto(user *models.User, path string) error
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
	GetUserProgress(course models.Course) (float32, error)
	GetUserProgressForAllCourses(studentID int) ([]repo.CourseWithProgressResponse, error)
	GetTeachersCourses(userID int) ([]models.Course, error)
	UpdateCourseImage(course models.Course) error

	AddCourse(course models.Course) (models.Course, error)

	GetTeachersCourse(userID, courseID int) (models.Course, error)
	//lessons
	GetLessons(studentID int, period string) ([]models.LessonResponse, error)
	GetLesson(lessonID int) (models.Lesson, error)
	GetLessonsByCourseID(courseID int) ([]models.Lesson, error)
	AddLesson(lesson models.Lesson) (models.Lesson, error)
	DeleteLesson(lessonId int) error

	UpdateLesson(lesson models.Lesson) (models.Lesson, error)

	//materials
	GetMaterials(lessonId int) ([]models.Material, error)
	GetMaterial(materialID int) (models.Material, error)

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
	RandomGenerateTest(testID int) (models.Test, error)
	CreateTest(test models.Test) (models.Test, error)
	DeleteTest(testID int) error

	//tasks
	GetTask(taskID int) (models.Task, error)
	AddStudentsTask(studentTask models.StudentTask) (models.StudentTask, error)
	UpdateStudentTask(studentTask models.StudentTask) error
	GetStudentTasks(taskID, studentID int) ([]models.StudentTask, error)
	GetStudentTask(taskID, studentID int) (models.StudentTask, error)
	GetFinalScore(taskID, studentID int) (float32, error)
	AddTask(task models.Task) (models.Task, error)
	DeleteTask(taskID int) error
	GetAllTasks() ([]models.StudentTask, error)
	GetTaskAnswersForTeacher(taskID int) ([]models.StudentTaskResponse, error)

	UpdateTask(task models.Task) error

	AddModule(module models.Module) error

	GetQuestionsForModule(moduleID int) ([]models.Question, error)
	AddQuestion(question models.Question) (models.Question, error)
	AddAnswer(answer models.AnswerVariant) (models.AnswerVariant, error)

	GetModule(moduleID int) (models.Module, error)
	DeleteModule(moduleID int) error
	UpdateModule(module models.Module) error
	UpdateTest(test models.Test) error

	CreateMaterial(material models.Material) error
	UpdateMaterial(material models.Material) error
	DeleteMaterial(materialID int) error

	CreateStudentTest(testID, studentID int) (models.CompletedTest, error)
	GetStudentsTests(testID int) ([]models.CompletedTestResponse, error)

	DeleteAnswer(answerID int) error
	DeleteQuestion(questionID int) error
	UpdateQuestion(question models.Question) error
	GetQuestion(questionID int) (models.Question, error)
	GetAnswerVariant(answerID int) (models.AnswerVariant, error)
	UpdateAnswerVariant(answerVariant models.AnswerVariant) error
}
