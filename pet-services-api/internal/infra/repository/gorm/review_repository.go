package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	reviewdom "pet-services-api/internal/domain/review"
	"pet-services-api/internal/models"
)

// ReviewRepository implementa review.Repository com GORM.
type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *reviewdom.Review) error {
	mr := toModelReview(rev)
	return r.db.WithContext(ctx).Create(mr).Error
}

func (r *ReviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*reviewdom.Review, error) {
	var m models.Review
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainReview(&m), nil
}

func (r *ReviewRepository) FindByRequestID(ctx context.Context, requestID uuid.UUID) (*reviewdom.Review, error) {
	var m models.Review
	if err := r.db.WithContext(ctx).First(&m, "request_id = ?", requestID).Error; err != nil {
		return nil, err
	}
	return toDomainReview(&m), nil
}

func (r *ReviewRepository) ExistsByRequestID(ctx context.Context, requestID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Review{}).Where("request_id = ?", requestID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ReviewRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, page, limit int) ([]*reviewdom.Review, int64, error) {
	base := r.db.WithContext(ctx).Model(&models.Review{}).
		Where("provider_id = ?", providerID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var ms []models.Review
	if err := applyPagination(base, page, limit).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*reviewdom.Review, 0, len(ms))
	for i := range ms {
		res = append(res, toDomainReview(&ms[i]))
	}
	return res, total, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id).Error
}
