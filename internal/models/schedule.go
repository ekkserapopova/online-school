package models

import "time"

type Schedule struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	StudentID int       `gorm:"unique;not null" json:"studentId"`
	Courses   []Course  `gorm:"many2many:schedules_courses;" json:"courses"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
