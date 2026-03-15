package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Breed struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null;index:idx_breed_species_name,unique" json:"name"`
	SpeciesID     string         `gorm:"type:char(26);not null;index:idx_breed_species_name,unique;index" json:"species_id"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Species Species `gorm:"foreignKey:SpeciesID" json:"specie"`
}

func (Breed) TableName() string {
	return "breeds"
}

func (b *Breed) ToEntity() *entities.Breed {
	return &entities.Breed{
		Base: entities.Base{
			ID:            b.ID,
			Active:        b.Active,
			CreatedAt:     b.CreatedAt,
			UpdatedAt:     b.UpdatedAt,
			DeactivatedAt: b.DeactivatedAt,
		},
		Name:      b.Name,
		SpeciesID: b.SpeciesID,
	}
}

func (b *Breed) FromEntity(entity *entities.Breed) {
	b.ID = entity.ID
	b.Name = entity.Name
	b.SpeciesID = entity.SpeciesID
	b.Active = entity.Active
	b.CreatedAt = entity.CreatedAt
	b.UpdatedAt = entity.UpdatedAt
	b.DeactivatedAt = entity.DeactivatedAt
}
