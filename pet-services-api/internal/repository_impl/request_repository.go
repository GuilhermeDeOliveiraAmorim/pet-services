package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) entities.RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(request *entities.Request) error {
	var model models.Request
	model.FromEntity(request)
	return r.db.Create(&model).Error
}

func (r *requestRepository) ExistsPending(userID, serviceID, petID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Request{}).
		Where("user_id = ? AND service_id = ? AND pet_id = ? AND status = ?", userID, serviceID, petID, entities.RequestStatuses.Pending).
		Count(&count).Error
	return count > 0, err
}
