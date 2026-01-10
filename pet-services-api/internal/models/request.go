package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index" json:"user_id"`
	ProviderID    string         `gorm:"type:char(26);not null;index" json:"provider_id"`
	ServiceID     string         `gorm:"type:char(26);not null;index" json:"service_id"`
	PetID         string         `gorm:"type:char(26);not null;index" json:"pet_id"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	RejectReason  string         `gorm:"type:varchar(500)" json:"reject_reason"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider"`
	Service  Service  `gorm:"foreignKey:ServiceID" json:"service"`
	Pet      Pet      `gorm:"foreignKey:PetID" json:"pet"`
}

func (Request) TableName() string {
	return "requests"
}
