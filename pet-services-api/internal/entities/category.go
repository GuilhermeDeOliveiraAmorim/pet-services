package entities

import (
	"pet-services-api/internal/exceptions"
)

type CategoryRepository interface {
	Create(category *Category) error
	FindByName(name string) (*Category, error)
}

type Category struct {
	Base
	Name string `json:"name"`
}

func NewCategory(name string) (*Category, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome da categoria ausente",
			Detail: "O nome da categoria é obrigatório",
		}))
	} else if len(name) > 50 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome da categoria deve ter no máximo 50 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Category{
		Base: *NewBase(),
		Name: name,
	}, nil
}
