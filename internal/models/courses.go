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
	TeacherID   int       `gorm:"not null" json:"teacherID"`
	Lessons     []Lesson  `gorm:"foreignKey:CourseID" json:"lessons"`
	Tests       []Test    `gorm:"foreignKey:CourseID" json:"tests"`
	Payments    []Payment `gorm:"foreignKey:CourseID" json:"payments"`
}

type Lesson struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Start       time.Time  `gorm:"not null" json:"start"`
	End         time.Time  `gorm:"not null" json:"end"`
	CourseID    int        `gorm:"not null" json:"courseID"`
	Homeworks   []Homework `gorm:"foreignKey:LessonID" json:"homeworks"`
	Materials   []Material `gorm:"foreignKey:LessonID" json:"materials"`
}

type Material struct { //надо связь добавить
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	File     string `json:"file"`
	LessonID int    `json:"lessonID"`
}

type Homework struct {
	ID                   int       `gorm:"primaryKey" json:"id"`
	Deadline             time.Time `json:"deadline"`
	File                 string    `json:"file"`
	Result               int       `json:"result"`
	Comment              string    `json:"comment"`               //комментарий от проверяющего
	ImplementationStatus string    `json:"implementation_status"` //не выполнено/на проверке/выполнено
	LessonID             int       `gorm:"not null" json:"lessonID"`
	StudentID            int       `gorm:"not null" json:"studentID"`
	Student              Student   `gorm:"foreignKey:StudentID" json:"student"`
	Status               bool      `gorm:"not null" json:"status"` //опубликовано или нет
}

type Test struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	LessonID    int        `json:"lessonID"`
	Lesson      Lesson     `gorm:"foreignKey:LessonID" json:"lesson"`
	CourseID    int        `gorm:"not null" json:"courseID"`
	Questions   []Question `gorm:"foreignKey:TestID" json:"questions"`
}

type Question struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Text   string `gorm:"type:text" json:"text"`
	Answer string `json:"answer"`
	TestID int    `json:"test_id"`
}

type Answer struct {
	ID         int  `gorm:"primaryKey" json:"id"`
	QuestionID int  `gorm:"not null" json:"questionID"`
	StudentID  int  `gorm:"not null" json:"studentID"`
	Answer     int  `json:"answer"`
	Result     bool `json:"result"`
}
