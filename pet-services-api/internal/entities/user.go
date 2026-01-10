package entities

import "pet-services-api/internal/exceptions"

type User struct {
	Base
	Name          string  `json:"name"`
	Login         Login   `json:"login"`
	Phone         Phone   `json:"phone"`
	Address       Address `json:"address"`
	EmailVerified bool    `json:"email_verified"`
}

func NewUser(name string, login Login, phone Phone, address Address) (*User, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome do cliente ausente",
			Detail: "O nome do cliente é obrigatório",
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
		EmailVerified: false,
	}, nil
}
