package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	ProviderID    string         `gorm:"type:char(26);not null;index:idx_provider_active,priority:1" json:"provider_id"`
	Name          string         `gorm:"type:varchar(100);not null;index" json:"name"`
	Description   string         `gorm:"type:text;not null" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2)" json:"price"`
	PriceMinimum  float64        `gorm:"type:decimal(10,2);index:idx_price_range,priority:1" json:"price_minimum"`
	PriceMaximum  float64        `gorm:"type:decimal(10,2);index:idx_price_range,priority:2" json:"price_maximum"`
	Duration      int            `gorm:"type:int" json:"duration"`
	Active        bool           `gorm:"default:true;index;index:idx_provider_active,priority:2" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Provider   Provider   `gorm:"foreignKey:ProviderID" json:"provider"`
	Photos     []Photo    `gorm:"many2many:service_photos;constraint:OnDelete:CASCADE" json:"photos"`
	Categories []Category `gorm:"many2many:service_categories;constraint:OnDelete:CASCADE" json:"categories"`
	Tags       []Tag      `gorm:"many2many:service_tags;constraint:OnDelete:CASCADE" json:"tags"`
	Requests   []Request  `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE" json:"requests"`
}

func (Service) TableName() string {
	return "services"
}
