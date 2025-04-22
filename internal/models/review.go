package models

import "time"

type Review struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Text        string    `gorm:"type:text" json:"text"`
	Assessment  int       `gorm:"not null" json:"assessment"` //от 1 до 5
	StudentID   int       `gorm:"not null" json:"student_id"`
	Student     User      `gorm:"foreignKey:StudentID" json:"student"`
	CourseID    int       `gorm:"not null" json:"course_id"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course"`
	IsPublished bool      `gorm:"not null" json:"is_published"` //опубликован или нет
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
