package models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);uniqueIndex;not null" json:"user_id"`
	BusinessName  string         `gorm:"type:varchar(100);not null;index" json:"business_name"`
	Description   string         `gorm:"type:text" json:"description"`
	PriceRange    string         `gorm:"type:varchar(10);index:idx_search,priority:3" json:"price_range"`
	AverageRating float64        `gorm:"type:decimal(3,2);default:0;index:idx_search,priority:2" json:"average_rating"`
	Active        bool           `gorm:"default:true;index;index:idx_search,priority:1" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Street       string  `gorm:"type:varchar(255);not null" json:"street"`
	Number       string  `gorm:"type:varchar(20);not null" json:"number"`
	Neighborhood string  `gorm:"type:varchar(100);not null" json:"neighborhood"`
	City         string  `gorm:"type:varchar(100);not null;index:idx_location" json:"city"`
	ZipCode      string  `gorm:"type:varchar(20);not null" json:"zip_code"`
	State        string  `gorm:"type:varchar(2);not null;index:idx_location" json:"state"`
	Country      string  `gorm:"type:varchar(50);not null;default:'BR'" json:"country"`
	Complement   string  `gorm:"type:varchar(255)" json:"complement"`
	Latitude     float64 `gorm:"type:decimal(10,7);index:idx_geo,priority:1" json:"latitude"`
	Longitude    float64 `gorm:"type:decimal(10,7);index:idx_geo,priority:2" json:"longitude"`

	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Photos   []Photo   `gorm:"many2many:provider_photos;constraint:OnDelete:CASCADE" json:"photos"`
	Services []Service `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE" json:"services"`
	Reviews  []Review  `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE" json:"reviews"`
	Requests []Request `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE" json:"requests"`
}

func (Provider) TableName() string {
	return "providers"
}
