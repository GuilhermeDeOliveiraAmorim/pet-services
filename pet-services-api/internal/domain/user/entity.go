package user

import (
	"time"

	"github.com/google/uuid"
)

// UserType representa o tipo de usuário
type UserType string

const (
	UserTypeOwner    UserType = "owner"
	UserTypeProvider UserType = "provider"
)

// User representa um usuário no sistema
type User struct {
	ID            uuid.UUID
	Email         Email
	EmailVerified bool
	Password      string // hash bcrypt
	Name          string
	Phone         Phone
	Type          UserType
	Location      *Location
	DeletedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewUser cria um novo usuário
func NewUser(email, name, phone string, userType UserType) (*User, error) {
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	phoneVO, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     emailVO,
		Name:      name,
		Phone:     phoneVO,
		Type:      userType,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// SetPassword define a senha do usuário (deve ser hash bcrypt)
func (u *User) SetPassword(hashedPassword string) {
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
}

// SetLocation define a localização do usuário
func (u *User) SetLocation(latitude, longitude float64, address Address) {
	u.Location = &Location{
		Latitude:  latitude,
		Longitude: longitude,
		Address:   address,
	}
	u.UpdatedAt = time.Now()
}

// UpdateProfile atualiza informações do perfil
func (u *User) UpdateProfile(name, phone string) error {
	if name != "" {
		u.Name = name
	}

	if phone != "" {
		phoneVO, err := NewPhone(phone)
		if err != nil {
			return err
		}
		u.Phone = phoneVO
	}

	u.UpdatedAt = time.Now()
	return nil
}

// SoftDelete marca o usuário como deletado.
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}

// IsDeleted verifica se o usuário foi deletado.
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// VerifyEmail marca o email como verificado.
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}
