package models

import (
	"time"

	"github.com/google/uuid"
)

// RequestStatus representa o status da solicitação.
type RequestStatus string

const (
	RequestStatusPending   RequestStatus = "pending"
	RequestStatusAccepted  RequestStatus = "accepted"
	RequestStatusRejected  RequestStatus = "rejected"
	RequestStatusCompleted RequestStatus = "completed"
	RequestStatusCancelled RequestStatus = "cancelled"
)

// PetType representa o tipo de pet atendido.
type PetType string

const (
	PetTypeDog     PetType = "dog"
	PetTypeCat     PetType = "cat"
	PetTypeBird    PetType = "bird"
	PetTypeRabbit  PetType = "rabbit"
	PetTypeFish    PetType = "fish"
	PetTypeRodent  PetType = "rodent"
	PetTypeReptile PetType = "reptile"
	PetTypeOther   PetType = "other"
)

// ServiceRequest modela a solicitação persistida.
type ServiceRequest struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey"`
	OwnerID         uuid.UUID     `gorm:"type:uuid;not null;index:idx_requests_owner;index:idx_requests_owner_created,composite:idx_requests_owner_created;foreignKey:OwnerID;references:ID"`
	ProviderID      uuid.UUID     `gorm:"type:uuid;not null;index:idx_requests_provider;index:idx_requests_provider_created,composite:idx_requests_provider_created;foreignKey:ProviderID;references:ID"`
	ServiceType     string        `gorm:"size:120;not null"`
	PetName         string        `gorm:"size:120;not null"`
	PetType         PetType       `gorm:"type:varchar(32);not null;index:idx_requests_pet_type"`
	PetBreed        string        `gorm:"size:120"`
	PetAge          int           `gorm:"not null;default:0"`
	PetWeight       float64       `gorm:"type:numeric(6,2);default:0"`
	PetNotes        string        `gorm:"size:500"`
	PreferredDate   time.Time     `gorm:"not null;index:idx_requests_preferred_date"`
	PreferredTime   string        `gorm:"size:5"`
	AdditionalNotes string        `gorm:"size:1000"`
	Status          RequestStatus `gorm:"type:varchar(16);not null;index:idx_requests_status;index:idx_requests_status_created,composite:idx_requests_status_created"`
	RejectionReason string        `gorm:"size:500"`
	CreatedAt       time.Time     `gorm:"autoCreateTime;index:idx_requests_created_at;index:idx_requests_owner_created,composite:idx_requests_owner_created;index:idx_requests_provider_created,composite:idx_requests_provider_created;index:idx_requests_status_created,composite:idx_requests_status_created"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime"`

	Owner    *User `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:RESTRICT"`
	Provider *User `gorm:"foreignKey:ProviderID;references:ID;constraint:OnDelete:RESTRICT"`
}

// TableName define o nome da tabela no banco.
func (ServiceRequest) TableName() string {
	return "requests"
}
