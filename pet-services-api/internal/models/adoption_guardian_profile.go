package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type AdoptionGuardianProfile struct {
	ID             string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID         string         `gorm:"type:char(26);not null;uniqueIndex" json:"user_id"`
	DisplayName    string         `gorm:"type:varchar(120);not null" json:"display_name"`
	GuardianType   string         `gorm:"type:varchar(20);not null" json:"guardian_type"`
	Document       string         `gorm:"type:varchar(30)" json:"document"`
	Phone          string         `gorm:"type:varchar(30)" json:"phone"`
	Whatsapp       string         `gorm:"type:varchar(30)" json:"whatsapp"`
	About          string         `gorm:"type:text" json:"about"`
	CityID         string         `gorm:"type:varchar(50);index" json:"city_id"`
	StateID        string         `gorm:"type:varchar(10);index" json:"state_id"`
	ApprovalStatus string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"approval_status"`
	ApprovedBy     string         `gorm:"type:char(26)" json:"approved_by"`
	ApprovedAt     *time.Time     `json:"approved_at"`
	Active         bool           `gorm:"default:true;index" json:"active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt  *time.Time     `json:"deactivated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (AdoptionGuardianProfile) TableName() string {
	return "adoption_guardian_profiles"
}

func (m *AdoptionGuardianProfile) ToEntity() *entities.AdoptionGuardianProfile {
	return &entities.AdoptionGuardianProfile{
		Base: entities.Base{
			ID:            m.ID,
			Active:        m.Active,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeactivatedAt: m.DeactivatedAt,
		},
		UserID:         m.UserID,
		DisplayName:    m.DisplayName,
		GuardianType:   m.GuardianType,
		Document:       m.Document,
		Phone:          m.Phone,
		Whatsapp:       m.Whatsapp,
		About:          m.About,
		CityID:         m.CityID,
		StateID:        m.StateID,
		ApprovalStatus: m.ApprovalStatus,
		ApprovedBy:     m.ApprovedBy,
		ApprovedAt:     m.ApprovedAt,
	}
}

func (m *AdoptionGuardianProfile) FromEntity(e *entities.AdoptionGuardianProfile) {
	m.ID = e.ID
	m.UserID = e.UserID
	m.DisplayName = e.DisplayName
	m.GuardianType = e.GuardianType
	m.Document = e.Document
	m.Phone = e.Phone
	m.Whatsapp = e.Whatsapp
	m.About = e.About
	m.CityID = e.CityID
	m.StateID = e.StateID
	m.ApprovalStatus = e.ApprovalStatus
	m.ApprovedBy = e.ApprovedBy
	m.ApprovedAt = e.ApprovedAt
	m.Active = e.Active
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
	m.DeactivatedAt = e.DeactivatedAt
}
