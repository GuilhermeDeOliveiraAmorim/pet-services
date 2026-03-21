package models

import (
	"pet-services-api/internal/entities"
	"time"
)

type AdoptionApplication struct {
	ID              string     `gorm:"type:char(26);primaryKey" json:"id"`
	ListingID       string     `gorm:"size:26;index" json:"listing_id"`
	ApplicantUserID string     `gorm:"size:26;index" json:"applicant_user_id"`
	Status          string     `gorm:"type:varchar(50);index" json:"status"`
	Motivation      string     `gorm:"type:text" json:"motivation"`
	HousingType     string     `gorm:"type:varchar(100)" json:"housing_type"`
	HasOtherPets    bool       `json:"has_other_pets"`
	PetExperience   string     `gorm:"type:text" json:"pet_experience"`
	FamilyMembers   int        `json:"family_members"`
	AgreesHomeVisit bool       `json:"agrees_home_visit"`
	ContactPhone    string     `gorm:"type:varchar(30)" json:"contact_phone"`
	NotesInternal   string     `gorm:"type:text" json:"notes_internal"`
	ReviewedBy      string     `gorm:"size:26" json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	Active          bool       `gorm:"default:true;index" json:"active"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at"`
}

func (AdoptionApplication) TableName() string {
	return "adoption_applications"
}

func (m *AdoptionApplication) ToEntity() *entities.AdoptionApplication {
	return &entities.AdoptionApplication{
		Base: entities.Base{
			ID:            m.ID,
			Active:        m.Active,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeactivatedAt: m.DeletedAt,
		},
		ListingID:       m.ListingID,
		ApplicantUserID: m.ApplicantUserID,
		Status:          m.Status,
		Motivation:      m.Motivation,
		HousingType:     m.HousingType,
		HasOtherPets:    m.HasOtherPets,
		PetExperience:   m.PetExperience,
		FamilyMembers:   m.FamilyMembers,
		AgreesHomeVisit: m.AgreesHomeVisit,
		ContactPhone:    m.ContactPhone,
		NotesInternal:   m.NotesInternal,
		ReviewedBy:      m.ReviewedBy,
		ReviewedAt:      m.ReviewedAt,
	}
}

func (m *AdoptionApplication) FromEntity(e *entities.AdoptionApplication) {
	m.ID = e.ID
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
	m.DeletedAt = e.DeactivatedAt
	m.ListingID = e.ListingID
	m.ApplicantUserID = e.ApplicantUserID
	m.Status = e.Status
	m.Motivation = e.Motivation
	m.HousingType = e.HousingType
	m.HasOtherPets = e.HasOtherPets
	m.PetExperience = e.PetExperience
	m.FamilyMembers = e.FamilyMembers
	m.AgreesHomeVisit = e.AgreesHomeVisit
	m.ContactPhone = e.ContactPhone
	m.NotesInternal = e.NotesInternal
	m.ReviewedBy = e.ReviewedBy
	m.ReviewedAt = e.ReviewedAt
	m.Active = true
}
