package entities

import (
	"pet-services-api/internal/exceptions"
	"regexp"
)

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
	} else if !isNumeric(countryCode) {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Código do país inválido",
			Detail: "O código do país deve conter apenas números",
		}))
	} else if len(countryCode) < 1 || len(countryCode) > 3 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Código do país inválido",
			Detail: "O código do país deve ter entre 1 e 3 dígitos",
		}))
	}

	if areaCode == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O código de área é obrigatório",
			Detail: "O campo código de área é obrigatório",
		}))
	} else if !isNumeric(areaCode) {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Código de área inválido",
			Detail: "O código de área deve conter apenas números",
		}))
	} else if len(areaCode) < 2 || len(areaCode) > 3 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Código de área inválido",
			Detail: "O código de área deve ter entre 2 e 3 dígitos",
		}))
	}

	if number == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "O número é obrigatório",
			Detail: "O campo número é obrigatório",
		}))
	} else if !isNumeric(number) {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Número inválido",
			Detail: "O número deve conter apenas dígitos",
		}))
	} else if len(number) < 8 || len(number) > 9 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Número inválido",
			Detail: "O número deve ter entre 8 e 9 dígitos",
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

func isNumeric(s string) bool {
	matched, _ := regexp.MatchString("^[0-9]+$", s)
	return matched
}

func (p *Phone) String() string {
	return "+" + p.CountryCode + " (" + p.AreaCode + ") " + p.Number
}
