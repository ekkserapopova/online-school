package models

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Surname   string    `gorm:"not null" json:"surname"`
	Name      string    `gorm:"not null" json:"name"`
	Birth     time.Time `gorm:"type:date" json:"birth"`
	Photo     string    `json:"photo"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Phone     string    `gorm:"unique;not null" json:"phone"`
	Password  string    `gorm:"not null" json:"password"`
	Salt      string    `gorm:"not null" json:"salt"`
	IsTeacher bool      ` json:"is_teacher"`
	IsActive  bool      `gorm:"not null" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Overview  string    `gorm:"type:text" json:"overview"`
	//TODO: add to json
	Courses        []Course        `gorm:"many2many:students_courses;" json:"-"`
	CompletedTests []CompletedTest `gorm:"foreignKey:StudentID" json:"completed_tests"`
	StudentTasks   []StudentTask   `gorm:"foreignKey:StudentID" json:"student_tasks"`
	//Payments []Payment `gorm:"foreignKey:StudentID,constraint:fk_users_reviews" json:"-"`
}
