package models

import (
	"pet-services-api/internal/entities"
	"time"
)

type AdoptionApplicationEvent struct {
	ID            string     `gorm:"type:char(26);primaryKey" json:"id"`
	ApplicationID string     `gorm:"size:26;index" json:"application_id"`
	EventType     string     `gorm:"type:varchar(100)" json:"event_type"`
	ActorUserID   string     `gorm:"size:26" json:"actor_user_id"`
	PayloadJSON   string     `gorm:"type:jsonb" json:"payload_json"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at"`
}

func (AdoptionApplicationEvent) TableName() string {
	return "adoption_application_events"
}

func (m *AdoptionApplicationEvent) ToEntity() *entities.AdoptionApplicationEvent {
	return &entities.AdoptionApplicationEvent{
		Base: entities.Base{
			ID:            m.ID,
			Active:        true,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeactivatedAt: m.DeletedAt,
		},
		ApplicationID: m.ApplicationID,
		EventType:     m.EventType,
		ActorUserID:   m.ActorUserID,
		PayloadJSON:   m.PayloadJSON,
	}
}

func (m *AdoptionApplicationEvent) FromEntity(e *entities.AdoptionApplicationEvent) {
	m.ID = e.ID
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
	m.DeletedAt = e.DeactivatedAt
	m.ApplicationID = e.ApplicationID
	m.EventType = e.EventType
	m.ActorUserID = e.ActorUserID
	m.PayloadJSON = e.PayloadJSON
}
