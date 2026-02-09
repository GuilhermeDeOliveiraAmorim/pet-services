package repository_impl

import (
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

func (r *specieRepository) List() ([]*entities.Specie, error) {
	var modelsSpecies []models.Specie
	err := r.db.Order("name asc").Find(&modelsSpecies).Error
	if err != nil {
		return nil, err
	}

	species := make([]*entities.Specie, 0, len(modelsSpecies))
	for _, model := range modelsSpecies {
		species = append(species, model.ToEntity())
	}

	return species, nil
}
