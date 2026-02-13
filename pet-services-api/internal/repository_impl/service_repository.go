package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) entities.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *entities.Service) error {
	var model models.Service
	model.FromEntity(service)
	return r.db.Create(&model).Error
}
