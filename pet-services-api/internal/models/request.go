package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Request struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index:idx_user_status,priority:1" json:"user_id"`
	ProviderID    string         `gorm:"type:char(26);not null;index:idx_provider_status,priority:1" json:"provider_id"`
	ServiceID     string         `gorm:"type:char(26);not null;index" json:"service_id"`
	PetID         string         `gorm:"type:char(26);not null;index" json:"pet_id"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending';index;index:idx_user_status,priority:2;index:idx_provider_status,priority:2;index:idx_status_created,priority:1" json:"status"`
	RejectReason  string         `gorm:"type:varchar(500)" json:"reject_reason"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index:idx_status_created,priority:2" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	User     User     `gorm:"foreignKey:UserID" json:"user"`
	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider"`
	Service  Service  `gorm:"foreignKey:ServiceID" json:"service"`
	Pet      Pet      `gorm:"foreignKey:PetID" json:"pet"`
}

func (Request) TableName() string {
	return "requests"
}

func (r *Request) ToEntity() *entities.Request {
	return &entities.Request{
		Base: entities.Base{
			ID:            r.ID,
			Active:        r.Active,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
			DeactivatedAt: r.DeactivatedAt,
		},
		UserID:       r.UserID,
		ProviderID:   r.ProviderID,
		ServiceID:    r.ServiceID,
		Pet:          *r.Pet.ToEntity(),
		Notes:        r.Notes,
		Status:       r.Status,
		RejectReason: r.RejectReason,
	}
}

func (r *Request) FromEntity(entity *entities.Request) {
	r.ID = entity.ID
	r.UserID = entity.UserID
	r.ProviderID = entity.ProviderID
	r.ServiceID = entity.ServiceID
	r.PetID = entity.Pet.ID
	r.Notes = entity.Notes
	r.Status = entity.Status
	r.RejectReason = entity.RejectReason
	r.Active = entity.Active
	r.CreatedAt = entity.CreatedAt
	r.UpdatedAt = entity.UpdatedAt
	r.DeactivatedAt = entity.DeactivatedAt
}
