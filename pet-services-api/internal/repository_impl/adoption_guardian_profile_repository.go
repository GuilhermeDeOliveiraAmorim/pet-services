package repository_impl

import (
	"errors"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type adoptionGuardianProfileRepository struct {
	db *gorm.DB
}

func NewAdoptionGuardianProfileRepository(db *gorm.DB) entities.AdoptionGuardianProfileRepository {
	return &adoptionGuardianProfileRepository{db: db}
}

func (r *adoptionGuardianProfileRepository) Create(profile *entities.AdoptionGuardianProfile) error {
	var model models.AdoptionGuardianProfile
	model.FromEntity(profile)
	return r.db.Create(&model).Error
}

func (r *adoptionGuardianProfileRepository) FindByID(id string) (*entities.AdoptionGuardianProfile, error) {
	var model models.AdoptionGuardianProfile
	err := r.db.Where("id = ? AND active = ?", id, true).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.AdoptionGuardianProfileNotFoundError)
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *adoptionGuardianProfileRepository) FindByUserID(userID string) (*entities.AdoptionGuardianProfile, error) {
	var model models.AdoptionGuardianProfile
	err := r.db.Where("user_id = ? AND active = ?", userID, true).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(consts.AdoptionGuardianProfileNotFoundError)
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *adoptionGuardianProfileRepository) Update(profile *entities.AdoptionGuardianProfile) error {
	var model models.AdoptionGuardianProfile
	model.FromEntity(profile)
	return r.db.Save(&model).Error
}

func (r *adoptionGuardianProfileRepository) ListByApprovalStatus(status string, page, pageSize int) ([]*entities.AdoptionGuardianProfile, int64, error) {
	var profiles []models.AdoptionGuardianProfile
	var total int64

	query := r.db.Model(&models.AdoptionGuardianProfile{}).Where("active = ?", true)
	if status != "" {
		query = query.Where("approval_status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&profiles).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*entities.AdoptionGuardianProfile, len(profiles))
	for i, p := range profiles {
		result[i] = p.ToEntity()
	}
	return result, total, nil
}
