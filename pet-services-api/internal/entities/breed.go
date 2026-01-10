package entities

import "pet-services-api/internal/exceptions"

type Breed struct {
	Base
	Name     string `json:"name"`
	SpecieID string `json:"specie_id"`
}

func NewBreed(name string, specieID string) (*Breed, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome da raça ausente",
			Detail: "O nome da raça é obrigatório",
		}))
	} else if len(name) > 50 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome da raça deve ter no máximo 50 caracteres",
		}))
	}

	if specieID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da espécie ausente",
			Detail: "O ID da espécie é obrigatório",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Breed{
		Base:     *NewBase(),
		Name:     name,
		SpecieID: specieID,
	}, nil
}
