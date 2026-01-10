package entities

import "pet-services-api/internal/exceptions"

type Phone struct {
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}

func NewPhone(countryCode, areaCode, number string) (*Phone, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if countryCode == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O código do país é obrigatório",
			Detail: "O campo código do país é obrigatório",
		}))
	}

	if areaCode == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O código de área é obrigatório",
			Detail: "O campo código de área é obrigatório",
		}))
	}

	if number == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O número é obrigatório",
			Detail: "O campo número é obrigatório",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Phone{
		CountryCode: countryCode,
		AreaCode:    areaCode,
		Number:      number,
	}, nil
}

func (p *Phone) String() string {
	return "+" + p.CountryCode + " (" + p.AreaCode + ") " + p.Number
}
