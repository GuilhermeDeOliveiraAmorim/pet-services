package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Specie struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Breeds []Breed `gorm:"foreignKey:SpecieID;constraint:OnDelete:CASCADE" json:"breeds"`
}

func (Specie) TableName() string {
	return "species"
}

func (s *Specie) ToEntity() *entities.Specie {
	return &entities.Specie{
		Base: entities.Base{
			ID:            s.ID,
			Active:        s.Active,
			CreatedAt:     s.CreatedAt,
			UpdatedAt:     s.UpdatedAt,
			DeactivatedAt: s.DeactivatedAt,
		},
		Name: s.Name,
	}
}

func (s *Specie) FromEntity(entity *entities.Specie) {
	s.ID = entity.ID
	s.Name = entity.Name
	s.Active = entity.Active
	s.CreatedAt = entity.CreatedAt
	s.UpdatedAt = entity.UpdatedAt
	s.DeactivatedAt = entity.DeactivatedAt
}
