package entities

import (
	"fmt"
	"pet-services-api/internal/exceptions"
)

type Address struct {
	Street       string   `json:"street"`
	Number       string   `json:"number"`
	Neighborhood string   `json:"neighborhood"`
	City         string   `json:"city"`
	ZipCode      string   `json:"zip_code"`
	State        string   `json:"state"`
	Country      string   `json:"country"`
	Complement   string   `json:"complement"`
	Location     Location `json:"location"`
}

func NewAddress(street, number, neighborhood, city, zipCode, state, country, complement string, location Location) (*Address, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if street == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "A rua é obrigatória",
			Detail: "O campo rua é obrigatório",
		}))
	}

	if number == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O número é obrigatório",
			Detail: "O campo número é obrigatório",
		}))
	}

	if neighborhood == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O bairro é obrigatório",
			Detail: "O campo bairro é obrigatório",
		}))
	}

	if city == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "A cidade é obrigatória",
			Detail: "O campo cidade é obrigatório",
		}))
	}

	if zipCode == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O CEP é obrigatório",
			Detail: "O campo CEP é obrigatório",
		}))
	}

	if state == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O estado é obrigatório",
			Detail: "O campo estado é obrigatório",
		}))
	}

	if country == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O país é obrigatório",
			Detail: "O campo país é obrigatório",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Address{
		Street:       street,
		City:         city,
		State:        state,
		ZipCode:      zipCode,
		Country:      country,
		Number:       number,
		Neighborhood: neighborhood,
		Complement:   complement,
		Location:     location,
	}, nil
}

func (a *Address) FullAddress() string {
	return a.Street + ", " + a.Number + ", " + a.Neighborhood + ", " + a.City + ", " + a.State + ", " + a.Country + ", " + a.ZipCode + ", " + a.Complement + ", Location(" +
		fmt.Sprintf("Lat: %.6f, Lon: %.6f", a.Location.Latitude, a.Location.Longitude) + ")"
}
