package repo

import (
	"errors"
	"log"
	"onlineschool/internal/models"

	"gorm.io/gorm"
)

func (r *Repo) GetPayment(userID, courseID int) (models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Student").Preload("Course").Where("student_id = ? AND course_id = ?", userID, courseID).First(&payment).Error
	if err != nil {
		return payment, err
	}
	return payment, nil
}

func (r *Repo) AddPayment(userID, courseID int) (models.Payment, error) {
	var payment models.Payment

	payment.CourseID = courseID
	payment.StudentID = userID

	if err := r.db.Where("course_id = ? and student_id = ?", courseID, userID).First(&payment).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal("Strench error in getting payment")
			return payment, err // Произошла другая ошибка, не связанная с отсутствием записи
		}
	} else {
		log.Println("payment already exists")
		return payment, errors.New("payment already exists") // Платеж уже существует
	}

	log.Println("creating payment")

	var course models.Course

	err := r.db.Where("id = ?", courseID).First(&course).Error

	if err != nil {
		log.Fatal("Error of getting course for price")
		return payment, err
	}

	log.Println("Course for price successfull got")
	payment.Amount = course.Price

	if payment.Amount == 0 {
		payment.Status = "paid"
	} else {
		payment.Status = "not paid"
	}
	err = r.db.Create(&payment).Error

	// payment.Status = "not paid"
	// if payment.Amount == 0 {
	// 	err = r.EnrollStudent(userID, courseID)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// }

	if err != nil {
		log.Fatal("Error of creating payment")
		return payment, err
	}
	log.Println("Payment created")
	return payment, nil
}

func (r *Repo) UpdatePaymentsStatus(payment *models.Payment) error {
	var existingPayment models.Payment
	err := r.db.Where("id = ?", payment.ID).First(&existingPayment).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&existingPayment).Updates(payment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) HasStudentPaidForCourse(studentID, courseID int) (bool, error) {
	var count int64
	err := r.db.Table("payments").Where("student_id = ? AND course_id = ?", studentID, courseID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
