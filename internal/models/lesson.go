package models

import "time"

type Lesson struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Start       time.Time `gorm:"not null" json:"start"`
	End         time.Time `gorm:"not null" json:"end"`
	CourseID    int       `gorm:"not null" json:"courseID"`
	IsActive    bool      `gorm:"not null" json:"is_active"` //опубликован или нет
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	//Course      Course     `gorm:"foreignKey:CourseID" json:"course"`
	//TODO: add to json
	Homeworks []Homework `gorm:"foreignKey:LessonID" json:"-"`
	Materials []Material `gorm:"foreignKey:LessonID" json:"-"`
}

type Material struct { //надо связь добавить
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	File      string    `json:"file"`
	LessonID  int       `json:"lessonID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
