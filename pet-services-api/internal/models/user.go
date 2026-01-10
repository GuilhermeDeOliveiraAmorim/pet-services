package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:char(26);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	UserType      string         `gorm:"type:varchar(20);not null;index" json:"user_type"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"type:varchar(255);not null" json:"password"`
	CountryCode   string         `gorm:"type:varchar(3);not null" json:"country_code"`
	AreaCode      string         `gorm:"type:varchar(3);not null" json:"area_code"`
	PhoneNumber   string         `gorm:"type:varchar(15);not null" json:"phone_number"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	Active        bool           `gorm:"default:true;index" json:"active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Street       string  `gorm:"type:varchar(255);not null" json:"street"`
	Number       string  `gorm:"type:varchar(20);not null" json:"number"`
	Neighborhood string  `gorm:"type:varchar(100);not null" json:"neighborhood"`
	City         string  `gorm:"type:varchar(100);not null;index" json:"city"`
	ZipCode      string  `gorm:"type:varchar(20);not null" json:"zip_code"`
	State        string  `gorm:"type:varchar(2);not null;index" json:"state"`
	Country      string  `gorm:"type:varchar(50);not null;default:'BR'" json:"country"`
	Complement   string  `gorm:"type:varchar(255)" json:"complement"`
	Latitude     float64 `gorm:"type:decimal(10,7);index" json:"latitude"`
	Longitude    float64 `gorm:"type:decimal(10,7);index" json:"longitude"`

	Photos   []Photo   `gorm:"many2many:user_photos;constraint:OnDelete:CASCADE" json:"photos"`
	Pets     []Pet     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"pets"`
	Provider *Provider `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"provider,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"reviews"`
	Requests []Request `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"requests"`
}

func (User) TableName() string {
	return "users"
}
