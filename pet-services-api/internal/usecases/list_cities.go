package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/reference"
)

type ListCitiesInput struct {
	StateID *int `json:"state_id"`
}

type ListCitiesOutput struct {
	Cities []reference.City `json:"cities"`
}

type ListCitiesUseCase struct {
	service reference.Service
	logger  logging.LoggerInterface
}

func NewListCitiesUseCase(service reference.Service, logger logging.LoggerInterface) *ListCitiesUseCase {
	return &ListCitiesUseCase{service: service, logger: logger}
}

func (uc *ListCitiesUseCase) Execute(ctx context.Context, input ListCitiesInput) (*ListCitiesOutput, []exceptions.ProblemDetails) {
	if input.StateID != nil && *input.StateID <= 0 {
		return nil, uc.logger.LogBadRequest(ctx, "ListCitiesUseCase.Execute", "state_id inválido", errors.New("state_id deve ser maior que zero"))
	}

	cities, err := uc.service.ListCities(ctx, input.StateID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, "ListCitiesUseCase.Execute", "Erro ao buscar cidades", err)
	}

	uc.logger.LogInfo(ctx, "ListCitiesUseCase.Execute", "Listagem de cidades executada")
	return &ListCitiesOutput{Cities: cities}, nil
}