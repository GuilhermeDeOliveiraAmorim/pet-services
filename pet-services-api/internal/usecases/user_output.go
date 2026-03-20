package usecases

import (
	"time"

	"pet-services-api/internal/entities"
)

type UserLoginOutput struct {
	Email string `json:"email"`
}

type UserOutput struct {
	ID              string           `json:"id"`
	Active          bool             `json:"active"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	DeactivatedAt   *time.Time       `json:"deactivated_at"`
	Name            string           `json:"name"`
	UserType        string           `json:"user_type"`
	Login           UserLoginOutput  `json:"login"`
	Phone           entities.Phone   `json:"phone"`
	Address         entities.Address `json:"address"`
	EmailVerified   bool             `json:"email_verified"`
	ProfileComplete bool             `json:"profile_complete"`
	Photos          []entities.Photo `json:"photos"`
	Pets            []entities.Pet   `json:"pets"`
}

func NewUserOutput(user *entities.User) *UserOutput {
	if user == nil {
		return nil
	}

	photos := make([]entities.Photo, len(user.Photos))
	copy(photos, user.Photos)

	pets := make([]entities.Pet, len(user.Pets))
	copy(pets, user.Pets)

	return &UserOutput{
		ID:              user.ID,
		Active:          user.Active,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeactivatedAt:   user.DeactivatedAt,
		Name:            user.Name,
		UserType:        user.UserType,
		Login:           UserLoginOutput{Email: user.Login.Email},
		Phone:           user.Phone,
		Address:         user.Address,
		EmailVerified:   user.EmailVerified,
		ProfileComplete: user.ProfileComplete,
		Photos:          photos,
		Pets:            pets,
	}
}
