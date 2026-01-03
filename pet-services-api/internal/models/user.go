package models

import (
	"time"

	"github.com/google/uuid"
)

// UserType define os tipos de usuário suportados.
type UserType string

const (
	UserTypeOwner    UserType = "owner"
	UserTypeProvider UserType = "provider"
)

// User representa o usuário persistido.
type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email         string    `gorm:"size:254;not null;uniqueIndex:idx_users_email"`
	EmailVerified bool      `gorm:"default:false;index:idx_users_email_verified"`
	Password      string    `gorm:"size:255;not null"`
	Name          string    `gorm:"size:120;not null"`
	Phone         string    `gorm:"size:20;not null"`
	Type          UserType  `gorm:"type:varchar(16);not null;index:idx_users_type"`

	Latitude  *float64 `gorm:"type:decimal(10,7);column:latitude;index:idx_users_location"`
	Longitude *float64 `gorm:"type:decimal(10,7);column:longitude;index:idx_users_location"`

	Street     string `gorm:"size:150"`
	Number     string `gorm:"size:30"`
	Complement string `gorm:"size:150"`
	District   string `gorm:"size:120"`
	City       string `gorm:"size:120;index:idx_users_city"`
	State      string `gorm:"size:60;index:idx_users_state"`
	ZipCode    string `gorm:"size:20;index:idx_users_zip"`
	Country    string `gorm:"size:80;default:Brasil"`

	DeletedAt *time.Time `gorm:"index:idx_users_deleted_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_users_created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela no banco.
func (User) TableName() string {
	return "users"
}
