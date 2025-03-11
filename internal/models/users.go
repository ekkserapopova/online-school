package models

import (
	"time"
)

type User struct {
	ID      int       `gorm:"primaryKey" json:"id"`
	Surname string    `gorm:"not null" json:"surname"`
	Name    string    `gorm:"not null" json:"name"`
	Birth   time.Time `gorm:"type:date" json:"birth"`
	Photo   string    `json:"photo"`
	Email   string    `gorm:"unique; not null" json:"email"`
	Phone   string    `gorm:"unique; not null" json:"phone"`
}

type Teacher struct {
	User
	Overview string   `gorm:"type:text" json:"overview"`
	Courses  []Course `gorm:"foreignKey:TeacherID" json:"courses"`
}

type Student struct {
	User
	Payments []Payment `gorm:"foreignKey:StudentID" json:"payments"`
	Schedule Schedule  `gorm:"foreignKey:StudentID" json:"schedule"`
}

type Review struct {
	ID         int     `gorm:"primaryKey" json:"id"`
	Text       string  `gorm:"type:text" json:"text"`
	Assessment int     `gorm:"not null" json:"assessment"` //от 1 до 5
	StudentID  int     `gorm:"not null" json:"studentID"`
	Student    Student `gorm:"foreignKey:StudentID" json:"student"`
}

type Payment struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `gorm:"not null" json:"status"` //не оплачено/оплачено/в обработке
	Date      time.Time `json:"date"`
	StudentID int       `gorm:"not null" json:"studentID"`
	CourseID  int       `gorm:"not null" json:"courseID"`
}
