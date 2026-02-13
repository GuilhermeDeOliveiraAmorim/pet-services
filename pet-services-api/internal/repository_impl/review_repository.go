package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) entities.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *entities.Review) error {
	var model models.Review
	model.FromEntity(review)
	return r.db.Create(&model).Error
}

func (r *reviewRepository) List(providerID, userID string, page, pageSize int) ([]*entities.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.Model(&models.Review{}).
		Where("active = ?", true)

	if providerID != "" {
		query = query.Where("provider_id = ?", providerID)
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	entitiesList := make([]*entities.Review, len(reviews))
	for i, rev := range reviews {
		entitiesList[i] = rev.ToEntity()
	}

	return entitiesList, total, nil
}
