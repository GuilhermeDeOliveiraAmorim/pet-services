package models

import (
	"pet-services-api/internal/entities"
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
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Street       string  `gorm:"type:varchar(255);not null" json:"street"`
	Number       string  `gorm:"type:varchar(20);not null" json:"number"`
	Neighborhood string  `gorm:"type:varchar(100);not null" json:"neighborhood"`
	City         string  `gorm:"type:varchar(100);not null;index:idx_location" json:"city"`
	ZipCode      string  `gorm:"type:varchar(20);not null" json:"zip_code"`
	State        string  `gorm:"type:varchar(2);not null;index:idx_location" json:"state"`
	Country      string  `gorm:"type:varchar(50);not null;default:'BR'" json:"country"`
	Complement   string  `gorm:"type:varchar(255)" json:"complement"`
	Latitude     float64 `gorm:"type:decimal(10,7);index:idx_geo,priority:1" json:"latitude"`
	Longitude    float64 `gorm:"type:decimal(10,7);index:idx_geo,priority:2" json:"longitude"`

	Photos   []Photo   `gorm:"many2many:user_photos;constraint:OnDelete:CASCADE" json:"photos"`
	Pets     []Pet     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"pets"`
	Provider *Provider `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"provider,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"reviews"`
	Requests []Request `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"requests"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToEntity() *entities.User {
	var photos []entities.Photo
	for _, photo := range u.Photos {
		photos = append(photos, *photo.ToEntity())
	}

	var pets []entities.Pet
	for _, pet := range u.Pets {
		pets = append(pets, *pet.ToEntity())
	}

	return &entities.User{
		Base: entities.Base{
			ID:            u.ID,
			Active:        u.Active,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
			DeactivatedAt: u.DeactivatedAt,
		},
		Name:     u.Name,
		UserType: u.UserType,
		Login: entities.Login{
			Email:    u.Email,
			Password: u.Password,
		},
		Phone: entities.Phone{
			CountryCode: u.CountryCode,
			AreaCode:    u.AreaCode,
			Number:      u.PhoneNumber,
		},
		Address: entities.Address{
			Street:       u.Street,
			Number:       u.Number,
			Neighborhood: u.Neighborhood,
			City:         u.City,
			ZipCode:      u.ZipCode,
			State:        u.State,
			Country:      u.Country,
			Complement:   u.Complement,
			Location: entities.Location{
				Latitude:  u.Latitude,
				Longitude: u.Longitude,
			},
		},
		EmailVerified: u.EmailVerified,
		Photos:        photos,
		Pets:          pets,
	}
}

func (u *User) FromEntity(entity *entities.User) {
	u.ID = entity.ID
	u.Name = entity.Name
	u.UserType = entity.UserType
	u.Email = entity.Login.Email
	u.Password = entity.Login.Password
	u.CountryCode = entity.Phone.CountryCode
	u.AreaCode = entity.Phone.AreaCode
	u.PhoneNumber = entity.Phone.Number
	u.EmailVerified = entity.EmailVerified
	u.Active = entity.Active
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
	u.DeactivatedAt = entity.DeactivatedAt
	u.Street = entity.Address.Street
	u.Number = entity.Address.Number
	u.Neighborhood = entity.Address.Neighborhood
	u.City = entity.Address.City
	u.ZipCode = entity.Address.ZipCode
	u.State = entity.Address.State
	u.Country = entity.Address.Country
	u.Complement = entity.Address.Complement
	u.Latitude = entity.Address.Location.Latitude
	u.Longitude = entity.Address.Location.Longitude
}
