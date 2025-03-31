package repo

import (
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

func (r *Repo) GetCourse(id int) (models.Course, error) {
	var course models.Course
	err := r.db.Where("id = ?", id).First(&course).Error
	if err != nil {
		return course, err
	}
	return course, nil
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
