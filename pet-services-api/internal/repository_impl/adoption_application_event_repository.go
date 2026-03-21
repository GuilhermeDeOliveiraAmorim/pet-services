package repository_impl

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/models"

	"gorm.io/gorm"
)

type adoptionApplicationEventRepository struct {
	db *gorm.DB
}

func NewAdoptionApplicationEventRepository(db *gorm.DB) entities.AdoptionApplicationEventRepository {
	return &adoptionApplicationEventRepository{db: db}
}

func (r *adoptionApplicationEventRepository) Create(event *entities.AdoptionApplicationEvent) error {
	var m models.AdoptionApplicationEvent
	m.FromEntity(event)
	if err := r.db.Create(&m).Error; err != nil {
		return err
	}
	*event = *m.ToEntity()
	return nil
}

func (r *adoptionApplicationEventRepository) ListByApplicationID(applicationID string) ([]*entities.AdoptionApplicationEvent, error) {
	var events []models.AdoptionApplicationEvent

	if err := r.db.Where("application_id = ?", applicationID).
		Order("created_at ASC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.AdoptionApplicationEvent, len(events))
	for i, m := range events {
		result[i] = m.ToEntity()
	}

	return result, nil
}
