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
