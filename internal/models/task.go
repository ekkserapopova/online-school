package models

import "time"

type Task struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"unique" json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `gorm:"type:date" json:"deadline"`
	ModuleID    int       `json:"module_id"` //может быть null
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentTask struct {
	ID               int     `gorm:"primaryKey;autoIncrement"`
	StudentID        int     `json:"student_id"`
	TaskID           int     `json:"task_id"`
	Code             string  `json:"code"`
	Requirements     string  `json:"requirements_score"`
	Implementation   string  `json:"implementation_score"`
	BoundaryHandling string  `json:"boundary_handling_score"`
	Optimization     string  `json:"optimization_score"`
	Score            float32 `json:"score"`

	FullAnswerByLLm string    `json:"full_answer_by_llm"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	//CommentByTeacher string    `json:"comment_by_teacher"`

}
