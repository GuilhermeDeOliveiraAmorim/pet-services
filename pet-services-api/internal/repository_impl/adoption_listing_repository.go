package repository_impl

import (
	"errors"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type adoptionListingRepository struct {
	db *gorm.DB
}

func NewAdoptionListingRepository(db *gorm.DB) entities.AdoptionListingRepository {
	return &adoptionListingRepository{db: db}
}

func (r *adoptionListingRepository) Create(listing *entities.AdoptionListing) error {
	var m models.AdoptionListing
	m.FromEntity(listing)
	return r.db.Create(&m).Error
}

func (r *adoptionListingRepository) FindByID(id string) (*entities.AdoptionListing, error) {
	var m models.AdoptionListing
	err := r.db.
		Preload("Pet").
		Preload("Pet.Species").
		Preload("Pet.Photos").
		Preload("GuardianProfile").
		Where("id = ? AND active = ?", id, true).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.AdoptionListingNotFoundError)
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *adoptionListingRepository) Update(listing *entities.AdoptionListing) error {
	var m models.AdoptionListing
	m.FromEntity(listing)
	return r.db.Save(&m).Error
}

func (r *adoptionListingRepository) ListPublic(filters entities.AdoptionListingFilters, page, pageSize int) ([]*entities.AdoptionListing, int64, error) {
	var listings []models.AdoptionListing
	var total int64

	query := r.db.Model(&models.AdoptionListing{}).Where("active = ?", true)

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	} else {
		query = query.Where("status = ?", entities.AdoptionListingStatuses.Published)
	}
	if filters.Sex != "" {
		query = query.Where("sex = ?", filters.Sex)
	}
	if filters.Size != "" {
		query = query.Where("size = ?", filters.Size)
	}
	if filters.AgeGroup != "" {
		query = query.Where("age_group = ?", filters.AgeGroup)
	}
	if filters.CityID != "" {
		query = query.Where("city_id = ?", filters.CityID)
	}
	if filters.StateID != "" {
		query = query.Where("state_id = ?", filters.StateID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("Pet").
		Preload("Pet.Species").
		Preload("Pet.Photos").
		Preload("GuardianProfile").
		Order("published_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&listings).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.AdoptionListing, len(listings))
	for i, l := range listings {
		result[i] = l.ToEntity()
	}
	return result, total, nil
}

func (r *adoptionListingRepository) ListByGuardianProfileID(guardianProfileID string, page, pageSize int) ([]*entities.AdoptionListing, int64, error) {
	var listings []models.AdoptionListing
	var total int64

	query := r.db.Model(&models.AdoptionListing{}).
		Where("guardian_profile_id = ? AND active = ?", guardianProfileID, true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("Pet").
		Preload("Pet.Species").
		Preload("Pet.Photos").
		Preload("GuardianProfile").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&listings).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.AdoptionListing, len(listings))
	for i, l := range listings {
		result[i] = l.ToEntity()
	}
	return result, total, nil
}
