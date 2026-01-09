package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	requestdom "pet-services-api/internal/domain/request"
	"pet-services-api/internal/models"
)

// RequestRepository implementa request.Repository com GORM.
type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(ctx context.Context, req *requestdom.ServiceRequest) error {
	mr := toModelServiceRequest(req)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(mr).Error
	})
}

func (r *RequestRepository) FindByID(ctx context.Context, id uuid.UUID) (*requestdom.ServiceRequest, error) {
	var m models.ServiceRequest
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainServiceRequest(&m)
}

func (r *RequestRepository) Update(ctx context.Context, req *requestdom.ServiceRequest) error {
	mr := toModelServiceRequest(req)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Save(mr).Error
	})
}

func (r *RequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&models.ServiceRequest{}, "id = ?", id).Error
	})
}

func (r *RequestRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID, page, limit int) ([]*requestdom.ServiceRequest, int64, error) {
	base := r.db.WithContext(ctx).Model(&models.ServiceRequest{}).
		Where("owner_id = ?", ownerID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var ms []models.ServiceRequest
	if err := applyPagination(base, page, limit).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*requestdom.ServiceRequest, 0, len(ms))
	for i := range ms {
		rreq, err := toDomainServiceRequest(&ms[i])
		if err != nil {
			return nil, 0, err
		}
		res = append(res, rreq)
	}
	return res, total, nil
}

func (r *RequestRepository) FindByProviderID(ctx context.Context, providerID uuid.UUID, page, limit int) ([]*requestdom.ServiceRequest, int64, error) {
	base := r.db.WithContext(ctx).Model(&models.ServiceRequest{}).
		Where("provider_id = ?", providerID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var ms []models.ServiceRequest
	if err := applyPagination(base, page, limit).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*requestdom.ServiceRequest, 0, len(ms))
	for i := range ms {
		rreq, err := toDomainServiceRequest(&ms[i])
		if err != nil {
			return nil, 0, err
		}
		res = append(res, rreq)
	}
	return res, total, nil
}

func (r *RequestRepository) FindByStatus(ctx context.Context, status requestdom.Status, page, limit int) ([]*requestdom.ServiceRequest, int64, error) {
	base := r.db.WithContext(ctx).Model(&models.ServiceRequest{}).
		Where("status = ?", status)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var ms []models.ServiceRequest
	if err := applyPagination(base, page, limit).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*requestdom.ServiceRequest, 0, len(ms))
	for i := range ms {
		rreq, err := toDomainServiceRequest(&ms[i])
		if err != nil {
			return nil, 0, err
		}
		res = append(res, rreq)
	}
	return res, total, nil
}
