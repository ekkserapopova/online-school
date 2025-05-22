package repo

import "onlineschool/internal/models"

func (r *Repo) GetTask(taskID int) (models.Task, error) {
	task := models.Task{}
	err := r.db.Where("id = ?", taskID).First(&task).Error
	return task, err
}

func (r *Repo) AddStudentsTask(studentTask models.StudentTask) (models.StudentTask, error) {
	studentTask.Status = "in progress"
	err := r.db.Create(&studentTask).Error
	return studentTask, err
}

func (r *Repo) GetStudentTask(taskID, studentID int) (models.StudentTask, error) {
	studentTask := models.StudentTask{}
	err := r.db.Where("task_id = ? and student_id = ?", taskID, studentID).Order("created_at DESC").First(&studentTask).Error
	return studentTask, err
}

func (r *Repo) GetStudentTasks(taskID, studentID int) ([]models.StudentTask, error) {
	studentTasks := []models.StudentTask{}
	err := r.db.
		Where("task_id = ? AND student_id = ?", taskID, studentID).
		Order("created_at desc").
		Find(&studentTasks).Error

	return studentTasks, err
}

func (r *Repo) UpdateStudentTask(studentTask models.StudentTask) error {
	err := r.db.Save(&studentTask).Error
	return err
}

func (r *Repo) GetFinalScore(taskID, studentID int) (float32, error) {
	var mark float32
	err := r.db.
		Table("student_tasks").
		Select("avg(score)").
		Where("student_id = ? and task_id = ?", studentID, taskID).
		Scan(&mark).Error

	return mark, err
}

func (r *Repo) AddTask(task models.Task) (models.Task, error) {
	err := r.db.Create(&task).Error
	return task, err
}

func (r *Repo) DeleteTask(taskID int) error {
	err := r.db.Where("id = ?", taskID).Delete(&models.Task{}).Error
	return err
}

func (r *Repo) UpdateTask(task models.Task) error {
	err := r.db.Save(&task).Error
	return err
}

func (r *Repo) GetAllTasks() ([]models.StudentTask, error) {
	tasks := []models.StudentTask{}
	err := r.db.Find(&tasks).Error

	return tasks, err
}

func (r *Repo) GetTaskAnswersForTeacher(taskID int) ([]models.StudentTaskResponse, error) {
	studentTasks := []models.StudentTask{}
	err := r.db.Where("task_id = ?", taskID).Find(&studentTasks).Error
	var tasksResult []models.StudentTaskResponse
	for _, task := range studentTasks {
		var student models.User
		err = r.db.Where("id = ?", task.StudentID).First(&student).Error
		tasksResult = append(tasksResult, models.StudentTaskResponse{
			StudentTask:    task,
			StudentName:    student.Name,
			StudentSurname: student.Surname,
		})

	}
	return tasksResult, err
}
