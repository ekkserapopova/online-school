package models

import "time"

type Task struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"unique" json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `gorm:"type:date" json:"deadline"`
	IsActive    bool      `json:"is_active"`
	ModuleID    int       `json:"module_id"` //может быть null
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	StudentTasks []StudentTask `gorm:"foreignKey:TaskID" json:"student_tasks"`
}

type StudentTask struct {
	ID                    int     `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID             int     `json:"student_id"`
	TaskID                int     `json:"task_id"`
	Code                  string  `json:"code"`
	Requirements          string  `json:"requirements"`
	RequirementsScore     int     `json:"requirements_score"`
	Implementation        string  `json:"implementation"`
	ImplementationScore   int     `json:"implementation_score"`
	BoundaryHandling      string  `json:"boundary_handling"`
	BoundaryHandlingScore int     `json:"boundary_handling_score"`
	Optimization          string  `json:"optimization"`
	OptimizationScore     int     `json:"optimization_score"`
	Recommendation        string  `json:"recommendation"`
	Score                 float32 `json:"score"`
	Status                string  `json:"status"`

	CodeWithCommentsByLLM string    `json:"code_with_comments_by_llm"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	//CommentByTeacher string    `json:"comment_by_teacher"`

}

type StudentTaskResponse struct {
	StudentTask
	StudentName    string `json:"student_name"`
	StudentSurname string `json:"student_surname"`
}
