package entities

import "pet-services-api/internal/exceptions"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLocation(latitude, longitude float64) (*Location, []exceptions.ProblemDetails) {
	if latitude < -90 || latitude > 90 {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Latitude inválida",
				Detail: "A latitude deve estar entre -90 e 90 graus",
			}),
		}
	}

	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
