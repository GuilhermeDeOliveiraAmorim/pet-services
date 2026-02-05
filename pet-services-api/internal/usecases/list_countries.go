package usecases

import (
	"context"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/reference"
)

type ListCountriesOutput struct {
	Countries []reference.Country `json:"countries"`
}

type ListCountriesUseCase struct {
	service reference.Service
	logger  logging.LoggerInterface
}

func NewListCountriesUseCase(service reference.Service, logger logging.LoggerInterface) *ListCountriesUseCase {
	return &ListCountriesUseCase{service: service, logger: logger}
}

func (uc *ListCountriesUseCase) Execute(ctx context.Context) (*ListCountriesOutput, []exceptions.ProblemDetails) {
	uc.logger.LogInfo(ctx, "ListCountriesUseCase.Execute", "Listagem de países executada")
	return &ListCountriesOutput{Countries: uc.service.ListCountries()}, nil
}