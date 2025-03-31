package repo

import "onlineschool/internal/models"

func (r *Repo) GetMaterials(lessonID int) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.Where("lesson_id = ?", lessonID).Find(&materials).Error
	if err != nil {
		return materials, err
	}

	return materials, nil
}
