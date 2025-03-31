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
	IsActive    bool      `gorm:"not null" json:"is_active"` //активен или нет
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Teacher     User      `gorm:"foreignKey:TeacherID" json:"teacher"`
	Lessons     []Lesson  `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`          // добавлен omitempty
	Tests       []Test    `gorm:"foreignKey:CourseID" json:"tests,omitempty"`            // добавлен omitempty
	Payments    []Payment `gorm:"foreignKey:CourseID" json:"payments,omitempty"`         // добавлен omitempty
	Students    []User    `gorm:"many2many:students_courses;" json:"students,omitempty"` // добавлен omitempty
	// Languages []Language `gorm:"foreignKey:CourseID" json:"languages"`
}

type Language struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Courses   []Course  `gorm:"many2many:language_courses;" json:"courses"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Lesson struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Start       time.Time  `gorm:"not null" json:"start"`
	End         time.Time  `gorm:"not null" json:"end"`
	CourseID    int        `gorm:"not null" json:"courseID"`
	IsActive    bool       `gorm:"not null" json:"is_active"` //опубликован или нет
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Course      Course     `gorm:"foreignKey:CourseID" json:"course"`
	Homeworks   []Homework `gorm:"foreignKey:LessonID" json:"homeworks"`
	Materials   []Material `gorm:"foreignKey:LessonID" json:"materials"`
}

type Material struct { //надо связь добавить
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	File      string    `json:"file"`
	LessonID  int       `json:"lessonID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Student              User      `gorm:"foreignKey:StudentID" json:"student"`
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
	IsActive    bool       `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Question struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text" json:"text"`
	Answer    string    `json:"answer"`
	TestID    int       `json:"test_id"`
	Answers   []Answer  `gorm:"foreignKey:QuestionID" json:"answers"`
	Points    int       `json:"points"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Answer struct {
	ID            int       `gorm:"primaryKey" json:"id"`
	QuestionID    int       `gorm:"not null" json:"questionID"`
	StudentID     int       `gorm:"not null" json:"studentID"`
	StudentAnswer string    `json:"student_answer"`
	Result        bool      `json:"result"`
	PointsEarned  int       `json:"points_earned"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
