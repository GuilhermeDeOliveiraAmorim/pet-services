package models

import (
	"pet-services-api/internal/entities"
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
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Provider   Provider   `gorm:"foreignKey:ProviderID" json:"provider"`
	Photos     []Photo    `gorm:"many2many:service_photos;constraint:OnDelete:CASCADE" json:"photos"`
	Categories []Category `gorm:"many2many:service_categories;constraint:OnDelete:CASCADE" json:"categories"`
	Tags       []Tag      `gorm:"many2many:service_tags;constraint:OnDelete:CASCADE" json:"tags"`
	Requests   []Request  `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE" json:"requests"`
}

func (Service) TableName() string {
	return "services"
}

func (s *Service) ToEntity() *entities.Service {
	var photos []entities.Photo
	for _, photo := range s.Photos {
		photos = append(photos, *photo.ToEntity())
	}

	var categories []entities.Category
	for _, category := range s.Categories {
		categories = append(categories, *category.ToEntity())
	}

	var tags []entities.Tag
	for _, tag := range s.Tags {
		tags = append(tags, *tag.ToEntity())
	}

	return &entities.Service{
		Base: entities.Base{
			ID:            s.ID,
			Active:        s.Active,
			CreatedAt:     s.CreatedAt,
			UpdatedAt:     s.UpdatedAt,
			DeactivatedAt: s.DeactivatedAt,
		},
		ProviderID:   s.ProviderID,
		Name:         s.Name,
		Description:  s.Description,
		Price:        s.Price,
		PriceMinimum: s.PriceMinimum,
		PriceMaximum: s.PriceMaximum,
		Duration:     s.Duration,
		Photos:       photos,
		Categories:   categories,
		Tags:         tags,
	}
}

func (s *Service) FromEntity(entity *entities.Service) {
	s.ID = entity.ID
	s.ProviderID = entity.ProviderID
	s.Name = entity.Name
	s.Description = entity.Description
	s.Price = entity.Price
	s.PriceMinimum = entity.PriceMinimum
	s.PriceMaximum = entity.PriceMaximum
	s.Duration = entity.Duration
	s.Active = entity.Active
	s.CreatedAt = entity.CreatedAt
	s.UpdatedAt = entity.UpdatedAt
	s.DeactivatedAt = entity.DeactivatedAt
}
