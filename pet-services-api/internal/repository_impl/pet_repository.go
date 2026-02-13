package repository_impl

import (
	"errors"
	"pet-services-api/internal/consts"
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

func (r *petRepository) FindByID(id string) (*entities.Pet, error) {
	var model models.Pet
	err := r.db.
		Preload("Photos").
		First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.PetNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}
