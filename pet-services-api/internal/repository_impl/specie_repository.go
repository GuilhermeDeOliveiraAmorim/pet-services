package repository_impl

import (
	"errors"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type specieRepository struct {
	db *gorm.DB
}

func NewSpecieRepository(db *gorm.DB) entities.SpecieRepository {
	return &specieRepository{db: db}
}

func (r *specieRepository) List() ([]*entities.Species, error) {
	var modelsSpecies []models.Species
	err := r.db.Order("name asc").Find(&modelsSpecies).Error
	if err != nil {
		return nil, err
	}

	species := make([]*entities.Species, 0, len(modelsSpecies))
	for _, model := range modelsSpecies {
		species = append(species, model.ToEntity())
	}

	return species, nil
}

func (r *specieRepository) FindByID(id string) (*entities.Species, error) {
	var model models.Species
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.SpecieNotFoundErr)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}
