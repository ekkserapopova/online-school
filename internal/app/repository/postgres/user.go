package repo

import (
	"onlineschool/internal/models"

	"gorm.io/gorm"
)

func (r *Repo) CreateUser(user *models.User) error {

	return r.db.Create(user).Error
}

func (r *Repo) FindByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err == nil {
		return &user, nil
	}

	return &user, gorm.ErrRecordNotFound
}

func (r *Repo) FindByPhone(phone string) error {
	var user models.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err == nil {
		return nil
	}

	return gorm.ErrRecordNotFound
}

func (r *Repo) GetByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err == nil {
		return &user, nil
	}

	return &user, gorm.ErrRecordNotFound

}

func (r *Repo) AddPhoto(user *models.User, path string) error {
	user.Photo = path
	err := r.db.Where("id = ?", user.ID).Save(&user).Error

	return err
}
