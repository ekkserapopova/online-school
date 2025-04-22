package models

import "time"

type Homework struct {
	ID                   int       `gorm:"primaryKey" json:"id"`
	Deadline             time.Time `json:"deadline"`
	File                 string    `json:"file"`
	Result               int       `json:"result"`
	Comment              string    `json:"comment"`               //комментарий от проверяющего
	ImplementationStatus string    `json:"implementation_status"` //не выполнено/на проверке/выполнено
	LessonID             int       `gorm:"not null" json:"lesson_id"`
	StudentID            int       `gorm:"not null" json:"student_id"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Student              User      `gorm:"foreignKey:StudentID" json:"student"`
	Status               bool      `gorm:"not null" json:"status"` //опубликовано или нет
}
