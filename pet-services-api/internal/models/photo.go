package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	URL           string         `gorm:"type:varchar(500);not null" json:"url"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Users     []User     `gorm:"many2many:user_photos" json:"users,omitempty"`
	Providers []Provider `gorm:"many2many:provider_photos" json:"providers,omitempty"`
	Pets      []Pet      `gorm:"many2many:pet_photos" json:"pets,omitempty"`
	Services  []Service  `gorm:"many2many:service_photos" json:"services,omitempty"`
}

func (Photo) TableName() string {
	return "photos"
}
