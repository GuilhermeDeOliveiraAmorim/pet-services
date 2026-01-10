package models

import (
	"time"

	"gorm.io/gorm"
)

type Breed struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(50);not null;index" json:"name"`
	SpecieID      string         `gorm:"type:char(26);not null;index" json:"specie_id"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Specie Specie `gorm:"foreignKey:SpecieID" json:"specie"`
}

func (Breed) TableName() string {
	return "breeds"
}
