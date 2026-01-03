package models

import (
	"time"

	"github.com/google/uuid"
)

// Review representa avaliações persistidas.
type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	RequestID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_reviews_request;foreignKey:RequestID;references:ID"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_provider;index:idx_reviews_provider_created,composite:idx_reviews_provider_created;foreignKey:ProviderID;references:ID"`
	OwnerID    uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_owner;foreignKey:OwnerID;references:ID"`
	Rating     int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment    string    `gorm:"size:1000"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index:idx_reviews_created_at;index:idx_reviews_provider_created,composite:idx_reviews_provider_created"`

	Request  *ServiceRequest `gorm:"foreignKey:RequestID;references:ID;constraint:OnDelete:RESTRICT"`
	Provider *User           `gorm:"foreignKey:ProviderID;references:ID;constraint:OnDelete:RESTRICT"`
	Owner    *User           `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:RESTRICT"`
}

// TableName define o nome da tabela no banco.
func (Review) TableName() string {
	return "reviews"
}
