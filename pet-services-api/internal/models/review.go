package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index" json:"user_id"`
	ProviderID    string         `gorm:"type:char(26);not null;index:idx_provider_reviews,priority:1" json:"provider_id"`
	Rating        float64        `gorm:"type:decimal(3,2);not null;index:idx_provider_reviews,priority:3" json:"rating"`
	Comment       string         `gorm:"type:varchar(500)" json:"comment"`
	Active        bool           `gorm:"default:true;index:idx_provider_reviews,priority:2" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_provider_reviews,priority:4" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider"`
}

func (Review) TableName() string {
	return "reviews"
}
