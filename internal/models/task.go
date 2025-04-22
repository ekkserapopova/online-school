package models

import "time"

type Task struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"unique" json:"name"`
	Description string    `json:"description"`
	ModuleID    int       `json:"module_id"` //может быть null
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
