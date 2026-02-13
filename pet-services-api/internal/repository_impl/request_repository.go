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

func (r *requestRepository) ExistsCompleted(userID, providerID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Request{}).
		Where("user_id = ? AND provider_id = ? AND status = ?", userID, providerID, entities.RequestStatuses.Completed).
		Count(&count).Error
	return count > 0, err
}

func (r *requestRepository) List(userID, providerID, status string, page, pageSize int) ([]*entities.Request, int64, error) {
	var requests []models.Request
	var total int64

	query := r.db.Model(&models.Request{}).
		Preload("User").
		Preload("Provider").
		Preload("Service").
		Preload("Pet.Species")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if providerID != "" {
		query = query.Where("provider_id = ?", providerID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entities.Request, len(requests))
	for i, req := range requests {
		entities[i] = req.ToEntity()
	}

	return entities, total, nil
}
