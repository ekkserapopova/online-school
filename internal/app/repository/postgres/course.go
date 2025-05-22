package repo

import (
	"gorm.io/gorm"
	"log"
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
	err := r.db.Table("students_courses").Preload("Modules").
		Select("courses.*").
		Joins("join courses on students_courses.course_id = courses.id").
		Where("students_courses.user_id = ?", userID).
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *Repo) GetTeachersCourses(userID int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Preload("Modules.Tests").
		Where("teacher_id = ? ", userID).
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
		Select("courses.*").
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC").Where("is_active = ?", true)
		}).
		Preload("Modules.Tests", func(db *gorm.DB) *gorm.DB {
			return db.Preload("CompletedTests", "student_id = ? AND status = ?", userID, "completed").Where("is_active = ?", true)
		}).
		Preload("Modules.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Preload("StudentTasks", "student_id = ? AND status = ?", userID, "completed").Where("is_active = ?", true)
		}).
		Joins("JOIN courses ON students_courses.course_id = courses.id").
		Where("students_courses.course_id = ? AND students_courses.user_id = ?", courseID, userID).
		Take(&course).Error

	if err != nil {
		return course, err
	}

	return course, nil
}

func (r *Repo) GetTeachersCourse(userID, courseID int) (models.Course, error) {
	var course models.Course
	err := r.db.Table("courses").
		Preload("Languages").
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Lessons.Materials", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Tests", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Where("id = ? AND teacher_id = ?", courseID, userID).
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

func (r *Repo) GetUserProgress(course models.Course) (float32, error) {
	var countTasksAndTests int
	var completedTasksAndTests int

	for _, module := range course.Modules {
		countTasksAndTests += len(module.Tests)
		countTasksAndTests += len(module.Tasks)
		for _, test := range module.Tests {
			countTasksAndTests++
			completedTasksAndTests += len(test.CompletedTests)
		}
		for _, task := range module.Tasks {
			countTasksAndTests++
			completedTasksAndTests += len(task.StudentTasks)
		}
	}

	if countTasksAndTests == 0 {
		return 0, nil
	}

	return float32(completedTasksAndTests) / float32(countTasksAndTests) * 100, nil
}

type CourseWithProgressResponse struct {
	models.Course
	Progress float32 `json:"progress"`
}

func (r *Repo) GetUserProgressForAllCourses(studentID int) ([]CourseWithProgressResponse, error) {
	var courses []models.Course
	coursesWithProgress := []CourseWithProgressResponse{}

	err := r.db.Table("students_courses").
		Select("courses.*").
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Modules.Lessons").
		Preload("Modules.Tests", func(db *gorm.DB) *gorm.DB {
			return db.Preload("CompletedTests", "student_id = ? AND status = ?", studentID, "completed")
		}).
		Preload("Modules.Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Preload("StudentTasks", "student_id = ? AND status = ?", studentID, "completed")
		}).
		Joins("JOIN courses ON students_courses.course_id = courses.id").
		Where("students_courses.user_id = ?", studentID).
		Find(&courses).Order("id ASC").Error

	if err != nil {
		return nil, err
	}

	for _, course := range courses {
		progress, err := r.GetUserProgress(course)
		log.Print(progress)
		if err != nil {
			return nil, err
		}

		coursesWithProgress = append(coursesWithProgress, CourseWithProgressResponse{course, progress})
		if course.ID == 1 {
			log.Print(CourseWithProgressResponse{course, progress})
		}
	}

	return coursesWithProgress, nil
}

func (r *Repo) AddCourse(course models.Course) (models.Course, error) {
	type LanguageCourse struct {
		Language_id int
		Course_id   int
	}
	err := r.db.Create(&course).Error
	if err != nil {
		return course, err
	}
	//
	for _, lang := range course.Languages {
		log.Println("LANGUAGE ID:", lang.ID, course.ID)
	}
	//
	//for _, language := range course.Languages {
	//
	//	newLang := LanguageCourse{
	//		language.ID,
	//		course.ID,
	//	}
	//	log.Println(newLang)
	//	err = r.db.Table("language_courses").Create(&newLang).Error
	//	if err != nil {
	//		return course, err
	//	}
	//}
	return course, err
}

func (r *Repo) UpdateCourseImage(course models.Course) error {
	err := r.db.Save(&course).Error
	return err
}
