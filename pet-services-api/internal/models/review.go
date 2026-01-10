package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index" json:"user_id"`
	ProviderID    string         `gorm:"type:char(26);not null;index" json:"provider_id"`
	Rating        float64        `gorm:"type:decimal(3,2);not null;index" json:"rating"`
	Comment       string         `gorm:"type:varchar(500)" json:"comment"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider"`
}

func (Review) TableName() string {
	return "reviews"
}
