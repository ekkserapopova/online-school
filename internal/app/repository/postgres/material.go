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

func (r *Repo) GetMaterial(materialID int) (models.Material, error) {
	var material models.Material
	err := r.db.Where("id = ?", materialID).Find(&material).Error
	return material, err
}

func (r *Repo) CreateMaterial(material models.Material) error {
	return r.db.Create(&material).Error
}

func (r *Repo) UpdateMaterial(material models.Material) error {
	return r.db.Save(&material).Error
}

func (r *Repo) DeleteMaterial(materialID int) error {
	return r.db.Delete(&models.Material{}, "id = ?", materialID).Error
}
