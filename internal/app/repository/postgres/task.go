package repo

import "onlineschool/internal/models"

func (r *Repo) GetTask(taskID int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.Where("id = ?", taskID).First(task).Error
	return task, err
}

func (r *Repo) AddStudentsTask(studentTask models.StudentTask) (models.StudentTask, error) {
	err := r.db.Create(&studentTask).Error
	return studentTask, err
}

func (r *Repo) UpdateStudentTask(studentTask models.StudentTask) error {
	err := r.db.Save(&studentTask).Error
	return err
}
