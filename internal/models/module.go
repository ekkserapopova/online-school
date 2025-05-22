package models

import "time"

type Module struct {
	ID          int       `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	CourseID    int       `gorm:"not null" json:"course_id"`
	OpenForAll  bool      `gorm:"not null" json:"open_for_all"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Tests     []Test     `gorm:"foreignKey:ModuleID" json:"tests"`
	Lessons   []Lesson   `gorm:"foreignKey:ModuleID" json:"lessons"`
	Tasks     []Task     `gorm:"foreignKey:ModuleID" json:"tasks"`
	Questions []Question `gorm:"foreignKey:ModuleID" json:"questions"`
}
