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
		Preload("Species").
		First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.PetNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *petRepository) ListByUser(userID string, page, pageSize int) ([]*entities.Pet, int64, error) {
	var pets []models.Pet
	var total int64

	query := r.db.Model(&models.Pet{}).
		Preload("Photos").
		Preload("Species").
		Where("user_id = ? AND active = ?", userID, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&pets).Error; err != nil {
		return nil, 0, err
	}

	entitiesList := make([]*entities.Pet, len(pets))
	for i, pet := range pets {
		entitiesList[i] = pet.ToEntity()
	}

	return entitiesList, total, nil
}

func (r *petRepository) Update(pet *entities.Pet) error {
	var model models.Pet
	model.FromEntity(pet)
	return r.db.Save(&model).Error
}
