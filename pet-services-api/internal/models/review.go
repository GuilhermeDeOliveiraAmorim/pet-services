package models

import (
	"time"

	"github.com/google/uuid"
)

// Review representa avaliações persistidas.
type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	RequestID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_reviews_request"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_provider"`
	OwnerID    uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_owner"`
	Rating     int       `gorm:"not null;check:rating_between_1_5"`
	Comment    string    `gorm:"size:1000"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName define o nome da tabela no banco.
func (Review) TableName() string {
	return "reviews"
}
