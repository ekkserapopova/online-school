package repo

import "onlineschool/internal/models"

func (r *Repo) GetTest(testID int) (models.Test, error) {
	var test models.Test
	err := r.db.Preload("Questions").Where("id = ?", testID).First(&test).Error
	if err != nil {
		return test, err
	}

	return test, nil
}

func (r *Repo) GetTests(courseID int) ([]models.Test, error) {
	var tests []models.Test
	err := r.db.Preload("Questions").Where("course_id = ?", courseID).Find(&tests).Error
	if err != nil {
		return tests, err
	}

	return tests, nil
}
