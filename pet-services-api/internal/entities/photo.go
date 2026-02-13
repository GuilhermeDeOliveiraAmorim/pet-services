package entities

import "pet-services-api/internal/exceptions"

type Photo struct {
	Base
	URL string `json:"url"`
}

type PhotoRepository interface {
	CreateAndAttachToPet(petID string, photo *Photo) error
	CreateAndAttachToService(serviceID string, photo *Photo) error
	CreateAndAttachToUser(userID string, photo *Photo) error
	ReplaceUserPhoto(userID string, photo *Photo) error
}

func NewPhoto(url string) (*Photo, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if url == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "URL da foto ausente",
			Detail: "A URL da foto é obrigatória",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Photo{
		Base: *NewBase(),
		URL:  url,
	}, nil
}
