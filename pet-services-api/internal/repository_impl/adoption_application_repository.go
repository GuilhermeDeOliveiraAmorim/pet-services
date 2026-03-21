package repository_impl

import (
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type adoptionApplicationRepository struct {
	db *gorm.DB
}

func NewAdoptionApplicationRepository(db *gorm.DB) entities.AdoptionApplicationRepository {
	return &adoptionApplicationRepository{db: db}
}

func (r *adoptionApplicationRepository) Create(application *entities.AdoptionApplication) error {
	var m models.AdoptionApplication
	m.FromEntity(application)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	*application = *m.ToEntity()
	return nil
}

func (r *adoptionApplicationRepository) FindByID(id string) (*entities.AdoptionApplication, error) {
	var m models.AdoptionApplication
	if err := r.db.Where("id = ? AND active = ?", id, true).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.AdoptionApplicationNotFoundError)
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *adoptionApplicationRepository) FindByListingIDAndApplicantUserID(listingID, applicantUserID string) (*entities.AdoptionApplication, error) {
	var m models.AdoptionApplication
	if err := r.db.Where("listing_id = ? AND applicant_user_id = ? AND active = ?", listingID, applicantUserID, true).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.AdoptionApplicationNotFoundError)
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *adoptionApplicationRepository) Update(application *entities.AdoptionApplication) error {
	var m models.AdoptionApplication
	m.FromEntity(application)
	if err := r.db.Model(&m).Where("id = ? AND active = ?", m.ID, true).Updates(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r *adoptionApplicationRepository) ListByApplicantUserID(applicantUserID string, page, pageSize int) ([]*entities.AdoptionApplication, int64, error) {
	var applications []models.AdoptionApplication
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.AdoptionApplication{}).Where("applicant_user_id = ? AND active = ?", applicantUserID, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("applicant_user_id = ? AND active = ?", applicantUserID, true).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.AdoptionApplication, len(applications))
	for i, m := range applications {
		result[i] = m.ToEntity()
	}

	return result, total, nil
}

func (r *adoptionApplicationRepository) ListByListingID(listingID string, page, pageSize int) ([]*entities.AdoptionApplication, int64, error) {
	var applications []models.AdoptionApplication
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.AdoptionApplication{}).Where("listing_id = ? AND active = ?", listingID, true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("listing_id = ? AND active = ?", listingID, true).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.AdoptionApplication, len(applications))
	for i, m := range applications {
		result[i] = m.ToEntity()
	}

	return result, total, nil
}
