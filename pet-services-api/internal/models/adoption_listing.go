package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type AdoptionListing struct {
	ID                     string         `gorm:"type:char(26);primaryKey" json:"id"`
	PetID                  string         `gorm:"type:char(26);not null;index" json:"pet_id"`
	GuardianProfileID      string         `gorm:"type:char(26);not null;index" json:"guardian_profile_id"`
	Title                  string         `gorm:"type:varchar(140);not null" json:"title"`
	Description            string         `gorm:"type:text;not null" json:"description"`
	AdoptionReason         string         `gorm:"type:varchar(1000)" json:"adoption_reason"`
	Status                 string         `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"`
	Sex                    string         `gorm:"type:varchar(10)" json:"sex"`
	Size                   string         `gorm:"type:varchar(10)" json:"size"`
	AgeGroup               string         `gorm:"type:varchar(10)" json:"age_group"`
	Vaccinated             bool           `gorm:"default:false" json:"vaccinated"`
	Neutered               bool           `gorm:"default:false" json:"neutered"`
	Dewormed               bool           `gorm:"default:false" json:"dewormed"`
	SpecialNeeds           bool           `gorm:"default:false" json:"special_needs"`
	GoodWithChildren       bool           `gorm:"default:false" json:"good_with_children"`
	GoodWithDogs           bool           `gorm:"default:false" json:"good_with_dogs"`
	GoodWithCats           bool           `gorm:"default:false" json:"good_with_cats"`
	RequiresHouseScreening bool           `gorm:"default:false" json:"requires_house_screening"`
	CityID                 string         `gorm:"type:varchar(50);index" json:"city_id"`
	StateID                string         `gorm:"type:varchar(10);index" json:"state_id"`
	Latitude               float64        `gorm:"type:decimal(10,7);default:0" json:"latitude"`
	Longitude              float64        `gorm:"type:decimal(10,7);default:0" json:"longitude"`
	PublishedAt            *time.Time     `json:"published_at"`
	AdoptedAt              *time.Time     `json:"adopted_at"`
	Active                 bool           `gorm:"default:true;index" json:"active"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt          *time.Time     `json:"deactivated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Pet             Pet                     `gorm:"foreignKey:PetID" json:"pet"`
	GuardianProfile AdoptionGuardianProfile `gorm:"foreignKey:GuardianProfileID" json:"guardian_profile"`
}

func (AdoptionListing) TableName() string {
	return "adoption_listings"
}

func (m *AdoptionListing) ToEntity() *entities.AdoptionListing {
	return &entities.AdoptionListing{
		Base: entities.Base{
			ID:            m.ID,
			Active:        m.Active,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeactivatedAt: m.DeactivatedAt,
		},
		PetID:                  m.PetID,
		Pet:                    *m.Pet.ToEntity(),
		GuardianProfileID:      m.GuardianProfileID,
		GuardianProfile:        *m.GuardianProfile.ToEntity(),
		Title:                  m.Title,
		Description:            m.Description,
		AdoptionReason:         m.AdoptionReason,
		Status:                 m.Status,
		Sex:                    m.Sex,
		Size:                   m.Size,
		AgeGroup:               m.AgeGroup,
		Vaccinated:             m.Vaccinated,
		Neutered:               m.Neutered,
		Dewormed:               m.Dewormed,
		SpecialNeeds:           m.SpecialNeeds,
		GoodWithChildren:       m.GoodWithChildren,
		GoodWithDogs:           m.GoodWithDogs,
		GoodWithCats:           m.GoodWithCats,
		RequiresHouseScreening: m.RequiresHouseScreening,
		CityID:                 m.CityID,
		StateID:                m.StateID,
		Latitude:               m.Latitude,
		Longitude:              m.Longitude,
		PublishedAt:            m.PublishedAt,
		AdoptedAt:              m.AdoptedAt,
	}
}

func (m *AdoptionListing) FromEntity(e *entities.AdoptionListing) {
	m.ID = e.ID
	m.PetID = e.PetID
	m.GuardianProfileID = e.GuardianProfileID
	m.Title = e.Title
	m.Description = e.Description
	m.AdoptionReason = e.AdoptionReason
	m.Status = e.Status
	m.Sex = e.Sex
	m.Size = e.Size
	m.AgeGroup = e.AgeGroup
	m.Vaccinated = e.Vaccinated
	m.Neutered = e.Neutered
	m.Dewormed = e.Dewormed
	m.SpecialNeeds = e.SpecialNeeds
	m.GoodWithChildren = e.GoodWithChildren
	m.GoodWithDogs = e.GoodWithDogs
	m.GoodWithCats = e.GoodWithCats
	m.RequiresHouseScreening = e.RequiresHouseScreening
	m.CityID = e.CityID
	m.StateID = e.StateID
	m.Latitude = e.Latitude
	m.Longitude = e.Longitude
	m.PublishedAt = e.PublishedAt
	m.AdoptedAt = e.AdoptedAt
	m.Active = e.Active
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
	m.DeactivatedAt = e.DeactivatedAt
}
