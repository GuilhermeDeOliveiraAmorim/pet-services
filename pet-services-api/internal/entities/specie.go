package entities

import "pet-services-api/internal/exceptions"

type Species struct {
	Base
	Name string `json:"name"`
}

type SpecieRepository interface {
	List() ([]*Species, error)
	FindByID(id string) (*Species, error)
}

func NewSpecie(name string) (*Species, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome da espécie ausente",
			Detail: "O nome da espécie é obrigatório",
		}))
	} else if len(name) > 50 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome da espécie deve ter no máximo 50 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Species{
		Base: *NewBase(),
		Name: name,
	}, nil
}
