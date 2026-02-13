package repository_impl

import (
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) entities.ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) Create(provider *entities.Provider) error {
	var model models.Provider
	model.FromEntity(provider)
	return r.db.Create(&model).Error
}

func (r *providerRepository) FindByUserID(userID string) (*entities.Provider, error) {
	var model models.Provider
	err := r.db.First(&model, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.ProviderNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}
