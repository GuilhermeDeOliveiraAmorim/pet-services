package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type breedRepository struct {
	db *gorm.DB
}

func NewBreedRepository(db *gorm.DB) entities.BreedRepository {
	return &breedRepository{db: db}
}

func (r *breedRepository) ListBySpecies(speciesID string) ([]*entities.Breed, error) {
	var modelBreeds []models.Breed
	err := r.db.
		Where("species_id = ? AND active = ?", speciesID, true).
		Order("name asc").
		Find(&modelBreeds).Error
	if err != nil {
		return nil, err
	}

	breeds := make([]*entities.Breed, 0, len(modelBreeds))
	for _, model := range modelBreeds {
		breeds = append(breeds, model.ToEntity())
	}

	return breeds, nil
}
