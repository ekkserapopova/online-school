package models

// что-то на сложном
type Schedule struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	StudentID int      `gorm:"unique;not null" json:"studentId"`
	Courses   []Course `gorm:"many2many:schedule_courses;" json:"courses"`
}
