package models

import (
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
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User   User    `gorm:"foreignKey:UserID" json:"user"`
	Specie Specie  `gorm:"foreignKey:SpecieID" json:"specie"`
	Breed  Breed   `gorm:"foreignKey:BreedID" json:"breed"`
	Photos []Photo `gorm:"many2many:pet_photos;constraint:OnDelete:CASCADE" json:"photos"`
}

func (Pet) TableName() string {
	return "pets"
}
