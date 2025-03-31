package repo

import (
	"log"
	"onlineschool/internal/models"
	"sort"
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
		&models.Language{},
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

// func (r *Repo) GetTeachersList(name string) ([]models.Teacher, error) {
// 	var teachers []models.Teacher
// 	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&teachers).Error
// 	if err != nil {
// 		return teachers, err
// 	}
// 	return teachers, nil
// }

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

func (r *Repo) GetMaterials(lessonID int) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.Where("lesson_id = ?", lessonID).Find(&materials).Error
	if err != nil {
		return materials, err
	}

	return materials, nil
}

// func (r *Repo) GetStudent(studentID int) (models.Student, error) {
// 	var student models.Student
// 	err := r.db.Where("id = ?", studentID).First(&student).Error
// 	if err != nil {
// 		return student, err
// 	}

// 	return student, nil
// }

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

func (r *Repo) CreateUser(user *models.User) error {

	return r.db.Create(user).Error
}

func (r *Repo) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err == nil {
		return &user, nil
	}

	return &user, gorm.ErrRecordNotFound
}

func (r *Repo) FindByPhone(phone string) error {
	var user models.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err == nil {
		return nil
	}

	return gorm.ErrRecordNotFound
}

func (r *Repo) GetByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err == nil {
		return &user, nil
	}

	return &user, gorm.ErrRecordNotFound

}

func (r *Repo) GetLanguages(name string) ([]models.Language, error) {
	var languages []models.Language
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&languages).Error
	if err != nil {
		return languages, err
	}
	return languages, nil
}

func (r *Repo) EnrollStudent(userID, courseID int) error {
	var schedule models.Schedule
	schedule.StudentID = userID
	log.Println(schedule.ID)
	err := r.db.Where("student_id = ?", userID).FirstOrCreate(&schedule).Error
	if err != nil {
		return err
	}
	err = r.db.Table("students_courses").Create(map[string]interface{}{
		"student_id": userID,
		"course_id":  courseID,
	}).Error
	if err != nil {
		return err
	}

	if err := r.db.Table("schedules_courses").Create(map[string]interface{}{
		"schedule_id": schedule.ID, "course_id": courseID}).Error; err != nil {
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
