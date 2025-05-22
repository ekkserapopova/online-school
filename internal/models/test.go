package models

import (
	"github.com/lib/pq"
	"time"
)

type CompletedTest struct {
	ID        int       `gorm:"primary_key" json:"id"`
	TestID    *int      `gorm:"not null" json:"test_id"`
	StudentID int       `gorm:"not null" json:"student_id"`
	Status    string    `gorm:"not null" json:"status"`
	Points    int       `gorm:"not null default=0" json:"points"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

const (
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusStarted   = "in progress"
)

var validStatuses = map[string]bool{
	StatusCompleted: true,
	StatusFailed:    true,
	StatusStarted:   true,
}

func (c CompletedTest) checkStatus() bool {
	return validStatuses[c.Status]
}

type Test struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	Deadline       time.Time `gorm:"not null;type:date" json:"deadline"`
	CountQuestions int       `gorm:"default=0" json:"count_questions"`
	TimeLimit      int       `json:"time_limit"`
	ModuleID       int       `gorm:"not null" json:"module_id"`
	IsActive       bool      `gorm:"not null" json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Questions      []Question      `gorm:"foreignKey:TestID" json:"questions"`
	StudentAnswers []StudentAnswer `gorm:"foreignKey:TestID" json:"student_answers"`
	CompletedTests []CompletedTest `gorm:"foreignKey:TestID" json:"completed_tests"`
}

type Question struct {
	ID             int             `gorm:"primaryKey" json:"id"`
	Text           string          `gorm:"type:text" json:"text"`
	TestID         *int            `json:"test_id"`
	ModuleID       int             `json:"module_id"`
	Answers        []AnswerVariant `gorm:"foreignKey:QuestionID" json:"answers"`
	StudentAnswers []StudentAnswer `gorm:"foreignKey:QuestionID" json:"students_answers"`
	Points         int             `json:"points"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type AnswerVariant struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	NumberID   int       `gorm:"not null" json:"number_id"`
	Text       string    `gorm:"type:text" json:"text"`
	IsRight    bool      `gorm:"not null" json:"is_right"`
	QuestionID int       `gorm:"not null" json:"question_id"`
	TestID     *int      `json:"test_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentAnswer struct {
	ID                int           `gorm:"primaryKey" json:"id"`
	QuestionID        int           `gorm:"not null" json:"question_id"`
	StudentID         int           `gorm:"not null" json:"student_id"`
	TestID            *int          `json:"test_id"`
	SelectedAnswerIDs pq.Int64Array `gorm:"type:integer[]" json:"selected_answer_ids"`

	//AnswerVariant AnswerVariant `gorm:"foreignKey:SelectedAnswerID" json:"-"`
	Result       bool      `json:"result"`
	PointsEarned int       `json:"points_earned"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type AnswerResponse struct {
	ID int `json:"id"` //answer's id
	AnswerVariant
}

type RightAnswer struct {
	ID         int  `json:"id"`
	QuestionID int  `json:"question_id"`
	TestID     *int `json:"test_id"`
	AnswerID   int  `json:"answer_id"`
}

type CompletedTestResponse struct {
	ID             int       `gorm:"primary_key" json:"id"`
	TestID         *int      `gorm:"not null" json:"test_id"`
	StudentID      int       `gorm:"not null" json:"student_id"`
	Status         string    `gorm:"not null" json:"status"`
	Points         int       `gorm:"not null default=0" json:"points"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	StudentName    string    `json:"student_name"`
	StudentSurname string    `json:"student_surname"`
}
