package models

import (
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
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Breeds []Breed `gorm:"foreignKey:SpecieID;constraint:OnDelete:CASCADE" json:"breeds"`
}

func (Specie) TableName() string {
	return "species"
}
