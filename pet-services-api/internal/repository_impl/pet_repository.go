package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type petRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) entities.PetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(pet *entities.Pet) error {
	var model models.Pet
	model.FromEntity(pet)
	return r.db.Create(&model).Error
}
