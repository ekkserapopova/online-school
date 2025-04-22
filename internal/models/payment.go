package models

import "time"

type Payment struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `gorm:"not null" json:"status"` //не оплачено/оплачено/в обработке
	Date      time.Time `json:"date"`
	StudentID int       `gorm:"not null" json:"studentID"`
	Student   User      `gorm:"foreignKey:StudentID" json:"student"`
	CourseID  int       `gorm:"not null" json:"course_id"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
