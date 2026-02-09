package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type Pet struct {
	Base
	UserID string  `json:"user_id"`
	Name   string  `json:"name"`
	Specie Specie  `json:"specie"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
	Notes  string  `json:"notes"`
	Photos []Photo `json:"photos"`
}

type PetRepository interface {
	Create(pet *Pet) error
	FindByID(id string) (*Pet, error)
}

func NewPet(userID string, name string, specie Specie, age int, weight float64, notes string) (*Pet, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if name == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome do pet ausente",
			Detail: "O nome do pet é obrigatório",
		}))
	} else if len(name) > 50 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome muito longo",
			Detail: "O nome do pet deve ter no máximo 50 caracteres",
		}))
	}

	if len(notes) > 500 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Observações muito longas",
			Detail: "As observações devem ter no máximo 500 caracteres",
		}))
	}

	if age < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Idade do pet inválida",
			Detail: "A idade do pet não pode ser negativa",
		}))
	}

	if weight < 0 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Peso do pet inválido",
			Detail: "O peso do pet não pode ser negativo",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Pet{
		Base:   *NewBase(),
		UserID: userID,
		Name:   name,
		Specie: specie,
		Age:    age,
		Weight: weight,
		Notes:  notes,
	}, nil
}

func (p *Pet) AddPhoto(photo Photo) {
	timeNow := time.Now()
	p.Photos = append(p.Photos, photo)
	p.UpdatedAt = &timeNow
}

func (p *Pet) RemovePhoto(photo Photo) {
	timeNow := time.Now()
	var updatedPhotos []Photo
	for _, pht := range p.Photos {
		if pht.ID != photo.ID {
			updatedPhotos = append(updatedPhotos, pht)
		}
	}
	p.Photos = updatedPhotos
	p.UpdatedAt = &timeNow
}
