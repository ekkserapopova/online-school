package models

import (
	"time"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Surname   string    `gorm:"not null" json:"surname"`
	Name      string    `gorm:"not null" json:"name"`
	Birth     time.Time `gorm:"type:date" json:"birth"`
	Photo     string    `json:"photo"`
	Email     string    `gorm:"unique;not null" json:"email"` // исправлено: uniqueIndex вместо unique;
	Phone     string    `gorm:"unique;not null" json:"phone"` // исправлено: uniqueIndex вместо unique;
	Password  string    `gorm:"not null" json:"password"`
	Salt      string    `gorm:"not null" json:"-"`
	IsAdmin   bool      `gorm:"not null" json:"is_admin"`
	IsActive  bool      `gorm:"not null" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Overview  string    `gorm:"type:text" json:"overview"`
	Courses   []Course  `gorm:"many2many:students_courses;" json:"courses,omitempty"` // добавлен omitempty
	Payments  []Payment `gorm:"foreignKey:StudentID" json:"payments,omitempty"`       // добавлен omitempty
	Schedule  Schedule  `gorm:"foreignKey:StudentID" json:"schedule"`
}

type Review struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Text        string    `gorm:"type:text" json:"text"`
	Assessment  int       `gorm:"not null" json:"assessment"` //от 1 до 5
	StudentID   int       `gorm:"not null" json:"studentID"`
	Student     User      `gorm:"foreignKey:StudentID" json:"student"`
	CourseID    int       `gorm:"not null" json:"courseID"`
	Course      Course    `gorm:"foreignKey:CourseID" json:"course"`
	IsPublished bool      `gorm:"not null" json:"is_published"` //опубликован или нет
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Payment struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `gorm:"not null" json:"status"` //не оплачено/оплачено/в обработке
	Date      time.Time `json:"date"`
	StudentID int       `gorm:"not null" json:"studentID"`
	CourseID  int       `gorm:"not null" json:"courseID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
