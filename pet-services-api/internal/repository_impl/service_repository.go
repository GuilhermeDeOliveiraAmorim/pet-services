package repository_impl

import (
	"errors"

	"pet-services-api/internal/consts"
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

func (r *serviceRepository) FindByID(id string) (*entities.Service, error) {
	var model models.Service
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.ServiceNotFoundError)
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *serviceRepository) HasTag(serviceID, tagID string) (bool, error) {
	var count int64
	err := r.db.Table("service_tags").
		Where("service_id = ? AND tag_id = ?", serviceID, tagID).
		Count(&count).Error
	return count > 0, err
}

func (r *serviceRepository) AddTag(serviceID, tagID string) error {
	service := models.Service{ID: serviceID}
	tag := models.Tag{ID: tagID}
	return r.db.Model(&service).Association("Tags").Append(&tag)
}

func (r *serviceRepository) HasCategory(serviceID, categoryID string) (bool, error) {
	var count int64
	err := r.db.Table("service_categories").
		Where("service_id = ? AND category_id = ?", serviceID, categoryID).
		Count(&count).Error
	return count > 0, err
}

func (r *serviceRepository) AddCategory(serviceID, categoryID string) error {
	service := models.Service{ID: serviceID}
	category := models.Category{ID: categoryID}
	return r.db.Model(&service).Association("Categories").Append(&category)
}
