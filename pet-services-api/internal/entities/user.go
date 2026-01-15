package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type UserTypeENUM struct {
	Owner    string `json:"owner"`
	Provider string `json:"provider"`
}

var UserTypes = UserTypeENUM{
	Owner:    "owner",
	Provider: "provider",
}

type User struct {
	Base
	Name          string  `json:"name"`
	UserType      string  `json:"user_type"`
	Login         Login   `json:"login"`
	Phone         Phone   `json:"phone"`
	Address       Address `json:"address"`
	EmailVerified bool    `json:"email_verified"`
	Photos        []Photo `json:"photos"`
	Pets          []Pet   `json:"pets"`
}

func NewUser(name string, userType string, login Login, phone Phone, address Address) (*User, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome do cliente ausente",
			Detail: "O nome do cliente é obrigatório",
		}))
	} else if len(name) > 100 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome do cliente deve ter no máximo 100 caracteres",
		}))
	}

	if userType == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Tipo de usuário ausente",
			Detail: "O tipo de usuário é obrigatório",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &User{
		Base:          *NewBase(),
		Name:          name,
		Login:         login,
		Phone:         phone,
		Address:       address,
		UserType:      userType,
		EmailVerified: false,
	}, nil
}

func (u *User) IsOwner() bool {
	return u.UserType == UserTypes.Owner
}

func (u *User) IsProvider() bool {
	return u.UserType == UserTypes.Provider
}

func (u *User) MarkEmailVerified() {
	timeNow := time.Now()
	u.EmailVerified = true
	u.UpdatedAt = &timeNow
}

func (u *User) AddPhoto(photo Photo) {
	timeNow := time.Now()
	u.Photos = append(u.Photos, photo)
	u.UpdatedAt = &timeNow
}

func (u *User) RemovePhoto(photo Photo) {
	timeNow := time.Now()
	var updatedPhotos []Photo
	for _, pht := range u.Photos {
		if pht.ID != photo.ID {
			updatedPhotos = append(updatedPhotos, pht)
		}
	}
	u.Photos = updatedPhotos
	u.UpdatedAt = &timeNow
}

func (u *User) AddPet(pet Pet) {
	timeNow := time.Now()
	u.Pets = append(u.Pets, pet)
	u.UpdatedAt = &timeNow
}

func (u *User) RemovePet(petID string) {
	for i, pet := range u.Pets {
		if pet.ID == petID {
			timeNow := time.Now()
			u.Pets = append(u.Pets[:i], u.Pets[i+1:]...)
			u.UpdatedAt = &timeNow
			return
		}
	}
}
