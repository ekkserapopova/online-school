package repo

import "onlineschool/internal/models"

func (r *Repo) GetTask(taskID int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.Where("id = ?", taskID).First(task).Error
	return task, err
}
