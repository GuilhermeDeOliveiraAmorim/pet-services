package gormrepo

import (
	"context"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	providerdom "github.com/guilherme/pet-services-api/internal/domain/provider"
	userdom "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/models"
)

// ProviderRepository implementa provider.Repository com GORM.
type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) Create(ctx context.Context, p *providerdom.Provider) error {
	mp, err := toModelProvider(p)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(mp).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ProviderRepository) FindByID(ctx context.Context, id uuid.UUID) (*providerdom.Provider, error) {
	var m models.Provider
	if err := r.db.WithContext(ctx).
		Preload("Services").
		Preload("Photos").
		Preload("WorkingHours").
		First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainProvider(&m)
}

func (r *ProviderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*providerdom.Provider, error) {
	var m models.Provider
	if err := r.db.WithContext(ctx).
		Preload("Services").
		Preload("Photos").
		Preload("WorkingHours").
		First(&m, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return toDomainProvider(&m)
}

func (r *ProviderRepository) Update(ctx context.Context, p *providerdom.Provider) error {
	mp, err := toModelProvider(p)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Replace services/photos/working hours to keep consistency.
		if err := tx.Where("provider_id = ?", p.ID).Delete(&models.ProviderService{}).Error; err != nil {
			return err
		}
		if err := tx.Where("provider_id = ?", p.ID).Delete(&models.ProviderPhoto{}).Error; err != nil {
			return err
		}
		if err := tx.Where("provider_id = ?", p.ID).Delete(&models.ProviderWorkingHour{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Provider{}).Where("id = ?", p.ID).
			Updates(map[string]any{
				"business_name":   mp.BusinessName,
				"description":     mp.Description,
				"street":          mp.Street,
				"number":          mp.Number,
				"complement":      mp.Complement,
				"district":        mp.District,
				"city":            mp.City,
				"state":           mp.State,
				"zip_code":        mp.ZipCode,
				"country":         mp.Country,
				"latitude":        mp.Latitude,
				"longitude":       mp.Longitude,
				"price_range":     mp.PriceRange,
				"avg_rating":      mp.AvgRating,
				"total_reviews":   mp.TotalReviews,
				"is_active":       mp.IsActive,
				"approval_status": mp.ApprovalStatus,
				"moderation_note": mp.ModerationNote,
				"updated_at":      mp.UpdatedAt,
			}).Error; err != nil {
			return err
		}
		if len(mp.Services) > 0 {
			if err := tx.Create(&mp.Services).Error; err != nil {
				return err
			}
		}
		if len(mp.Photos) > 0 {
			if err := tx.Create(&mp.Photos).Error; err != nil {
				return err
			}
		}
		if len(mp.WorkingHours) > 0 {
			if err := tx.Create(&mp.WorkingHours).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Provider{}, "id = ?", id).Error
}

func (r *ProviderRepository) List(ctx context.Context, page, limit int) ([]*providerdom.Provider, int64, error) {
	var ms []models.Provider
	tx := applyPagination(r.db.WithContext(ctx), page, limit)

	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Provider{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.Preload("Services").Preload("Photos").Preload("WorkingHours").Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	providers := make([]*providerdom.Provider, 0, len(ms))
	for i := range ms {
		p, err := toDomainProvider(&ms[i])
		if err != nil {
			return nil, 0, err
		}
		providers = append(providers, p)
	}
	return providers, total, nil
}

func (r *ProviderRepository) FindActiveByLocation(ctx context.Context, latitude, longitude, radiusKM float64, page, limit int) ([]*providerdom.Provider, int64, error) {
	// Simple bounding box filter to reduce scan; fine for MVP without PostGIS.
	degRadius := radiusKM / 111.0
	latMin := latitude - degRadius
	latMax := latitude + degRadius
	lonRadius := degRadius / math.Cos(latitude*math.Pi/180)
	lonMin := longitude - lonRadius
	lonMax := longitude + lonRadius

	base := r.db.WithContext(ctx).Model(&models.Provider{}).
		Where("is_active = ?", true).
		Where("latitude BETWEEN ? AND ?", latMin, latMax).
		Where("longitude BETWEEN ? AND ?", lonMin, lonMax)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var ms []models.Provider
	tx := applyPagination(base, page, limit)
	if err := tx.Preload("Services").Preload("Photos").Preload("WorkingHours").Order("avg_rating DESC").Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	providers := make([]*providerdom.Provider, 0, len(ms))
	for i := range ms {
		p, err := toDomainProvider(&ms[i])
		if err != nil {
			return nil, 0, err
		}
		providers = append(providers, p)
	}
	return providers, total, nil
}

func (r *ProviderRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Provider{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateRating incrementa rating agregado (uso otimista simples).
func (r *ProviderRepository) UpdateRating(ctx context.Context, providerID uuid.UUID, newRating float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p models.Provider
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, "id = ?", providerID).Error; err != nil {
			return err
		}
		total := float64(p.TotalReviews)*p.AvgRating + newRating
		newTotalReviews := p.TotalReviews + 1
		newAvg := total / float64(newTotalReviews)
		return tx.Model(&models.Provider{}).
			Where("id = ?", providerID).
			Updates(map[string]any{
				"total_reviews": newTotalReviews,
				"avg_rating":    newAvg,
			}).Error
	})
}

// helper to rebuild provider from user location.
func copyUserAddress(u *userdom.User) userdom.Address {
	if u == nil || u.Location == nil {
		return userdom.Address{}
	}
	return u.Location.Address
}
