package entities

import "pet-services-api/internal/exceptions"

type Tag struct {
	Base
	Name string `json:"name"`
}

func NewTag(name string) (*Tag, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome da tag ausente",
			Detail: "O nome da tag é obrigatório",
		}))
	} else if len(name) > 30 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome da tag deve ter no máximo 30 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Tag{
		Base: *NewBase(),
		Name: name,
	}, nil
}
