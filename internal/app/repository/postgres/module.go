package repo

import "onlineschool/internal/models"

func (r *Repo) AddModule(module models.Module) error {
	err := r.db.Create(&module).Error
	return err
}

// TODO: что это?
func (r *Repo) AddModuleToCourse(courseID int, module models.Module) (models.Module, error) {
	err := r.AddModule(module)
	if err != nil {
		return module, err
	}
	course, err := r.GetCourse(courseID)
	if err != nil {
		return module, err
	}

	course.Modules = append(course.Modules, module)
	err = r.db.Save(&course).Error
	return module, err
}

func (r *Repo) GetModule(moduleID int) (models.Module, error) {
	module := models.Module{}
	err := r.db.Where("id = ?", moduleID).First(&module).Error
	return module, err
}

func (r *Repo) DeleteModule(moduleID int) error {
	err := r.db.Where("id = ?", moduleID).Delete(&models.Module{}).Error
	return err
}

func (r *Repo) UpdateModule(module models.Module) error {
	err := r.db.Save(&module).Error
	return err
}
