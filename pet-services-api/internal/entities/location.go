package entities

import "pet-services-api/internal/exceptions"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLocation(latitude, longitude float64) (*Location, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if latitude < -90 || latitude > 90 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Latitude inválida",
			Detail: "A latitude deve estar entre -90 e 90 graus",
		}))
	}

	if longitude < -180 || longitude > 180 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Longitude inválida",
			Detail: "A longitude deve estar entre -180 e 180 graus",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
