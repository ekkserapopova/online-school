package models

import (
	"time"
)

type Course struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Difficulty  int       `gorm:"not null" json:"difficulty"`
	Price       int       `gorm:"not null;default:0" json:"price"`
	TeacherID   int       `gorm:"not null" json:"teacher_id"`
	IsActive    bool      `gorm:"not null" json:"is_active"` //активен или нет
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Teacher     *User     `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Modules     []Module  `gorm:"foreignKey:CourseID" json:"modules,omitempty"`
	//Tests       []Test    `gorm:"foreignKey:CourseID" json:"-"`
	Payments []Payment `gorm:"foreignKey:CourseID" json:"-"`
	Students []User    `gorm:"many2many:students_courses;" json:"-"`
	// Languages []Language `gorm:"foreignKey:CourseID" json:"languages"`
}

type Language struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Courses   []Course  `gorm:"many2many:language_courses;" json:"courses"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
