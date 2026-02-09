package models

import (
	"pet-services-api/internal/entities"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	UserID        string         `gorm:"type:char(26);not null;index" json:"user_id"`
	Name          string         `gorm:"type:varchar(50);not null" json:"name"`
	SpecieID      string         `gorm:"type:char(26);not null;index" json:"specie_id"`
	BreedID       string         `gorm:"type:char(26);not null;index" json:"breed_id"`
	Age           int            `gorm:"type:int" json:"age"`
	Weight        float64        `gorm:"type:decimal(6,2)" json:"weight"`
	Notes         string         `gorm:"type:varchar(500)" json:"notes"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	User   User    `gorm:"foreignKey:UserID" json:"user"`
	Specie Specie  `gorm:"foreignKey:SpecieID" json:"specie"`
	Photos []Photo `gorm:"many2many:pet_photos;constraint:OnDelete:CASCADE" json:"photos"`
}

func (Pet) TableName() string {
	return "pets"
}

func (p *Pet) ToEntity() *entities.Pet {
	var photos []entities.Photo
	for _, photo := range p.Photos {
		photos = append(photos, *photo.ToEntity())
	}

	return &entities.Pet{
		Base: entities.Base{
			ID:            p.ID,
			Active:        p.Active,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			DeactivatedAt: p.DeactivatedAt,
		},
		UserID: p.UserID,
		Name:   p.Name,
		Specie: *p.Specie.ToEntity(),
		Age:    p.Age,
		Weight: p.Weight,
		Notes:  p.Notes,
		Photos: photos,
	}
}

func (p *Pet) FromEntity(entity *entities.Pet) {
	p.ID = entity.ID
	p.UserID = entity.UserID
	p.Name = entity.Name
	p.SpecieID = entity.Specie.ID
	p.Age = entity.Age
	p.Weight = entity.Weight
	p.Notes = entity.Notes
	p.Active = entity.Active
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
	p.DeactivatedAt = entity.DeactivatedAt
}
