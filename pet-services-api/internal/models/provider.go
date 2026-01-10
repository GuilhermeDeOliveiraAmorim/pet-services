package models

import (
	"pet-services-api/internal/entities"
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
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

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

func (p *Provider) ToEntity() *entities.Provider {
	var photos []entities.Photo
	for _, photo := range p.Photos {
		photos = append(photos, *photo.ToEntity())
	}

	var reviews []entities.Review
	for _, review := range p.Reviews {
		reviews = append(reviews, *review.ToEntity())
	}

	var requests []entities.Request
	for _, request := range p.Requests {
		requests = append(requests, *request.ToEntity())
	}

	return &entities.Provider{
		Base: entities.Base{
			ID:            p.ID,
			Active:        p.Active,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			DeactivatedAt: p.DeactivatedAt,
		},
		UserID:       p.UserID,
		BusinessName: p.BusinessName,
		Address: entities.Address{
			Street:       p.Street,
			Number:       p.Number,
			Neighborhood: p.Neighborhood,
			City:         p.City,
			ZipCode:      p.ZipCode,
			State:        p.State,
			Country:      p.Country,
			Complement:   p.Complement,
			Location: entities.Location{
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
			},
		},
		Description:   p.Description,
		PriceRange:    p.PriceRange,
		AverageRating: p.AverageRating,
		Photos:        photos,
		Reviews:       reviews,
		Requests:      requests,
	}
}

func (p *Provider) FromEntity(entity *entities.Provider) {
	p.ID = entity.ID
	p.UserID = entity.UserID
	p.BusinessName = entity.BusinessName
	p.Description = entity.Description
	p.PriceRange = entity.PriceRange
	p.AverageRating = entity.AverageRating
	p.Active = entity.Active
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
	p.DeactivatedAt = entity.DeactivatedAt
	p.Street = entity.Address.Street
	p.Number = entity.Address.Number
	p.Neighborhood = entity.Address.Neighborhood
	p.City = entity.Address.City
	p.ZipCode = entity.Address.ZipCode
	p.State = entity.Address.State
	p.Country = entity.Address.Country
	p.Complement = entity.Address.Complement
	p.Latitude = entity.Address.Location.Latitude
	p.Longitude = entity.Address.Location.Longitude
}
